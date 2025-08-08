package env_test

import (
	"os"
	"os/user"
	"reflect"
	"testing"
	"time"

	"github.com/andreimerlescu/env"
)

// withEnv is a helper to temporarily set an environment variable for a test.
// It ensures the original state is restored after the test function completes.
func withEnv(t *testing.T, key, value string, testFunc func()) {
	t.Helper()
	originalValue, wasSet := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("withEnv: failed to set env var %s: %v", key, err)
	}

	testFunc()

	// Restore original state
	if wasSet {
		if err := os.Setenv(key, originalValue); err != nil {
			t.Fatalf("withEnv: failed to restore original env var %s: %v", key, err)
		}
	} else {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("withEnv: failed to unset env var %s after test: %v", key, err)
		}
	}
}

// TestMain runs before all tests to configure the package for a clean test environment.
func TestMain(m *testing.M) {
	// Disable printing errors during tests to keep test output clean.
	// Individual tests can re-enable it if they need to test logging.
	env.PrintErrors = false
	env.UseLogger = false
	env.AllowPanic = false
	env.UseMagic = false // Disable magic for predictable test runs

	// Run all tests
	os.Exit(m.Run())
}

func TestExistence(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		if env.Exists("THIS_VAR_SHOULD_NOT_EXIST") {
			t.Error("Exists() returned true for an unset variable")
		}
		withEnv(t, "TEST_EXISTS", "hello", func() {
			if !env.Exists("TEST_EXISTS") {
				t.Error("Exists() returned false for a set variable")
			}
		})
	})

	t.Run("MustExist", func(t *testing.T) {
		// Test success case (should not panic)
		withEnv(t, "MUST_EXIST_VAR", "i exist", func() {
			env.MustExist("MUST_EXIST_VAR")
		})

		// Test panic case
		t.Run("panics when var is not set", func(t *testing.T) {
			// Enable panic for this specific test
			originalAllowPanic := env.AllowPanic
			env.AllowPanic = true
			defer func() { env.AllowPanic = originalAllowPanic }()

			defer func() {
				if r := recover(); r == nil {
					t.Error("MustExist() did not panic for an unset variable")
				}
			}()
			env.MustExist("THIS_VAR_REALLY_MUST_EXIST")
		})
	})
}

func TestTruthy(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		isSet     bool
		wantTrue  bool
		wantFalse bool
	}{
		{"unset", "", false, false, true},
		{"set true", "true", true, true, false},
		{"set false", "false", true, false, true},
		{"set 1", "1", true, true, false}, // FIX: strconv.ParseBool recognizes "1" as true
		{"set 0", "0", true, false, true}, // FIX: strconv.ParseBool recognizes "0" as false
		{"set mixed case True", "True", true, true, false},
		{"set junk string", "junk", true, false, true}, // Invalid bool results in fallback (false)
	}

	for _, tc := range cases {
		runTest := func(t *testing.T) {
			gotTrue := env.IsTrue("BOOL_TEST")
			if gotTrue != tc.wantTrue {
				t.Errorf("IsTrue() with value '%s' (set: %v): got %v, want %v", tc.value, tc.isSet, gotTrue, tc.wantTrue)
			}
			gotFalse := env.IsFalse("BOOL_TEST")
			if gotFalse != tc.wantFalse {
				t.Errorf("IsFalse() with value '%s' (set: %v): got %v, want %v", tc.value, tc.isSet, gotFalse, tc.wantFalse)
			}
		}

		if tc.isSet {
			withEnv(t, "BOOL_TEST", tc.value, func() {
				t.Run(tc.name, runTest)
			})
		} else {
			t.Run(tc.name, runTest)
		}
	}

	t.Run("AreTrue", func(t *testing.T) {
		withEnv(t, "T1", "true", func() {
			withEnv(t, "T2", "true", func() {
				withEnv(t, "F1", "false", func() {
					if !env.AreTrue("T1", "T2") {
						t.Error("AreTrue returned false for all true variables")
					}
					if env.AreTrue("T1", "F1") {
						t.Error("AreTrue returned true when one variable was false")
					}
					if env.AreTrue("T1", "UNSET") {
						t.Error("AreTrue returned true when one variable was unset")
					}
				})
			})
		})
	})

	t.Run("AreFalse", func(t *testing.T) {
		withEnv(t, "F1", "false", func() {
			withEnv(t, "F2", "false", func() {
				withEnv(t, "T1", "true", func() {
					if !env.AreFalse("F1", "F2", "UNSET_VAR") {
						t.Error("AreFalse returned false for all false/unset variables")
					}
					if env.AreFalse("F1", "T1") {
						t.Error("AreFalse returned true when one variable was true")
					}
				})
			})
		})
	})
}

func TestTypeGetters(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		fallback := "default"
		if val := env.String("UNSET_STRING", fallback); val != fallback {
			t.Errorf("String unset failed: got %s, want %s", val, fallback)
		}
		withEnv(t, "SET_STRING", "hello", func() {
			want := "hello"
			if val := env.String("SET_STRING", fallback); val != want {
				t.Errorf("String set failed: got %s, want %s", val, want)
			}
		})
	})

	t.Run("Int", func(t *testing.T) {
		fallback := 99
		if val := env.Int("UNSET_INT", fallback); val != fallback {
			t.Errorf("Int unset failed: got %d, want %d", val, fallback)
		}
		withEnv(t, "SET_INT", "123", func() {
			want := 123
			if val := env.Int("SET_INT", fallback); val != want {
				t.Errorf("Int set failed: got %d, want %d", val, want)
			}
		})
		withEnv(t, "INVALID_INT", "abc", func() {
			if val := env.Int("INVALID_INT", fallback); val != fallback {
				t.Errorf("Int invalid failed: got %d, want %d", val, fallback)
			}
		})
	})

	t.Run("Int64", func(t *testing.T) {
		fallback := int64(99)
		want := int64(1234567890)
		withEnv(t, "SET_INT64", "1234567890", func() {
			if val := env.Int64("SET_INT64", fallback); val != want {
				t.Errorf("Int64 set failed: got %d, want %d", val, want)
			}
		})
	})

	t.Run("Float64", func(t *testing.T) {
		fallback := 99.9
		want := 123.45
		withEnv(t, "SET_FLOAT64", "123.45", func() {
			if val := env.Float64("SET_FLOAT64", fallback); val != want {
				t.Errorf("Float64 set failed: got %f, want %f", val, want)
			}
		})
	})

	t.Run("Float32", func(t *testing.T) {
		fallback := float32(99.9)
		want := float32(123.45)
		withEnv(t, "SET_FLOAT32", "123.45", func() {
			if val := env.Float32("SET_FLOAT32", fallback); val != want {
				t.Errorf("Float32 set failed: got %f, want %f", val, want)
			}
		})
	})

	t.Run("Duration", func(t *testing.T) {
		fallback := 1 * time.Minute
		if val := env.Duration("UNSET_DUR", fallback); val != fallback {
			t.Errorf("Duration unset failed: got %v, want %v", val, fallback)
		}
		// Test parsing a duration string
		withEnv(t, "SET_DUR_STR", "10s", func() {
			want := 10 * time.Second
			if val := env.Duration("SET_DUR_STR", fallback); val != want {
				t.Errorf("Duration from string failed: got %v, want %v", val, want)
			}
		})
		// Test parsing an integer (as nanoseconds)
		withEnv(t, "SET_DUR_INT", "5000000000", func() {
			want := 5 * time.Second
			if val := env.Duration("SET_DUR_INT", fallback); val != want {
				t.Errorf("Duration from int failed: got %v, want %v", val, want)
			}
		})
	})

	t.Run("UnitDuration", func(t *testing.T) {
		fallback := 3 * time.Second
		unit := time.Second
		// Test unset
		if val := env.UnitDuration("UNSET_UDUR", 3, unit); val != fallback {
			t.Errorf("UnitDuration unset failed: got %v, want %v", val, fallback)
		}
		// Test parsing an integer to be multiplied by unit
		withEnv(t, "SET_UDUR_INT", "15", func() {
			want := 15 * time.Second
			if val := env.UnitDuration("SET_UDUR_INT", 3, unit); val != want {
				t.Errorf("UnitDuration from int failed: got %v, want %v", val, want)
			}
		})
		// Test parsing a full duration string (unit should be ignored)
		withEnv(t, "SET_UDUR_STR", "2m", func() {
			want := 2 * time.Minute
			if val := env.UnitDuration("SET_UDUR_STR", 3, unit); val != want {
				t.Errorf("UnitDuration from string failed: got %v, want %v", val, want)
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		fallback := []string{"default"}
		if !reflect.DeepEqual(env.List("UNSET_LIST", fallback), fallback) {
			t.Error("List unset failed")
		}
		// Test with default separator (,)
		withEnv(t, "SET_LIST", " a, b ,c ", func() {
			want := []string{"a", "b", "c"}
			got := env.List("SET_LIST", fallback)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("List default separator failed: got %v, want %v", got, want)
			}
		})
		// Test with custom separator
		originalSep := env.ListSeparator
		env.ListSeparator = "|"
		defer func() { env.ListSeparator = originalSep }()
		withEnv(t, "SET_LIST_CUSTOM", " x | y | z ", func() {
			want := []string{"x", "y", "z"}
			got := env.List("SET_LIST_CUSTOM", fallback)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("List custom separator failed: got %v, want %v", got, want)
			}
		})
	})

	t.Run("Map", func(t *testing.T) {
		fallback := map[string]string{"d": "e"}
		if !reflect.DeepEqual(env.Map("UNSET_MAP", fallback), fallback) {
			t.Error("Map unset failed")
		}
		withEnv(t, "SET_MAP", " k1=v1 , k2= v2 ", func() {
			want := map[string]string{"k1": "v1", "k2": "v2"}
			got := env.Map("SET_MAP", fallback)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Map default separator failed: got %v, want %v", got, want)
			}
		})
		// Test with custom separators
		originalMapSep := env.MapSeparator
		originalItemSep := env.MapItemSeparator
		env.MapSeparator = "&"
		env.MapItemSeparator = ":"
		defer func() {
			env.MapSeparator = originalMapSep
			env.MapItemSeparator = originalItemSep
		}()
		withEnv(t, "SET_MAP_CUSTOM", " a:1 & b: 2 ", func() {
			want := map[string]string{"a": "1", "b": "2"}
			got := env.Map("SET_MAP_CUSTOM", fallback)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Map custom separator failed: got %v, want %v", got, want)
			}
		})
	})
}

func TestDataValidation(t *testing.T) {
	listFallback := []string{}
	mapFallback := map[string]string{}

	t.Run("ListLength", func(t *testing.T) {
		withEnv(t, "LIST_LEN", "a,b,c", func() {
			if env.ListLength("LIST_LEN", listFallback) != 3 {
				t.Error("ListLength failed")
			}
		})
	})

	t.Run("ListIsLength", func(t *testing.T) {
		withEnv(t, "LIST_IS_LEN", "a,b", func() {
			if !env.ListIsLength("LIST_IS_LEN", listFallback, 2) {
				t.Error("ListIsLength returned false for correct length")
			}
			if env.ListIsLength("LIST_IS_LEN", listFallback, 3) {
				t.Error("ListIsLength returned true for incorrect length")
			}
		})
	})

	t.Run("ListContains", func(t *testing.T) {
		withEnv(t, "LIST_CONTAINS", "apple,Banana,ORANGE", func() {
			if !env.ListContains("LIST_CONTAINS", listFallback, "banana") {
				t.Error("ListContains returned false for an existing item (case-insensitive)")
			}
			if env.ListContains("LIST_CONTAINS", listFallback, "grape") {
				t.Error("ListContains returned true for a non-existent item")
			}
		})
	})

	t.Run("MapHasKey", func(t *testing.T) {
		withEnv(t, "MAP_HAS_KEY", "k1=v1,k2=v2", func() {
			if !env.MapHasKey("MAP_HAS_KEY", mapFallback, "k1") {
				t.Error("MapHasKey returned false for an existing key")
			}
			if env.MapHasKey("MAP_HAS_KEY", mapFallback, "k3") {
				t.Error("MapHasKey returned true for a non-existent key")
			}
		})
	})

	t.Run("MapHasKeys", func(t *testing.T) {
		withEnv(t, "MAP_HAS_KEYS", "k1=v1,k2=v2,k3=v3", func() {
			if !env.MapHasKeys("MAP_HAS_KEYS", mapFallback, "k1", "k3") {
				t.Error("MapHasKeys returned false when all keys exist")
			}
			if env.MapHasKeys("MAP_HAS_KEYS", mapFallback, "k1", "k4") {
				t.Error("MapHasKeys returned true when a key was missing")
			}
		})
	})

	t.Run("IntInRange", func(t *testing.T) {
		withEnv(t, "IN_RANGE", "50", func() {
			if !env.IntInRange("IN_RANGE", 0, 10, 100) {
				t.Error("IntInRange returned false for a value within range")
			}
			if env.IntInRange("IN_RANGE", 0, 60, 100) {
				t.Error("IntInRange returned true for a value outside range")
			}
			if !env.IntInRange("IN_RANGE", 0, 50, 100) {
				t.Error("IntInRange returned false for a value on the lower boundary")
			}
		})
	})

	t.Run("Int64LessThan", func(t *testing.T) {
		withEnv(t, "LESS_THAN", "9", func() {
			if !env.Int64LessThan("LESS_THAN", 0, 10) {
				t.Error("Int64LessThan returned false for a valid comparison")
			}
			if env.Int64LessThan("LESS_THAN", 0, 9) {
				t.Error("Int64LessThan returned true for an invalid comparison")
			}
		})
	})

	t.Run("Int64GreaterThan", func(t *testing.T) {
		withEnv(t, "GREATER_THAN_64", "100", func() {
			if !env.Int64GreaterThan("GREATER_THAN_64", 0, 99) {
				t.Error("Int64GreaterThan(100 > 99) returned false, want true")
			}
			if env.Int64GreaterThan("GREATER_THAN_64", 0, 100) {
				t.Error("Int64GreaterThan(100 > 100) returned true, want false")
			}
		})
	})

	t.Run("Int64InRange", func(t *testing.T) {
		withEnv(t, "IN_RANGE_64", "50", func() {
			if !env.Int64InRange("IN_RANGE_64", 0, 10, 100) {
				t.Error("Int64InRange(50 in [10, 100]) returned false, want true")
			}
			if env.Int64InRange("IN_RANGE_64", 0, 60, 100) {
				t.Error("Int64InRange(50 in [60, 100]) returned true, want false")
			}
			if !env.Int64InRange("IN_RANGE_64", 0, 50, 50) {
				t.Error("Int64InRange(50 in [50, 50]) returned false, want true")
			}
		})
	})

	t.Run("IntLessThan", func(t *testing.T) {
		withEnv(t, "LESS_THAN_INT", "9", func() {
			if !env.IntLessThan("LESS_THAN_INT", 0, 10) {
				t.Error("IntLessThan(9 < 10) returned false, want true")
			}
			if env.IntLessThan("LESS_THAN_INT", 0, 9) {
				t.Error("IntLessThan(9 < 9) returned true, want false")
			}
		})
	})

	t.Run("IntGreaterThan", func(t *testing.T) {
		withEnv(t, "GREATER_THAN_INT", "100", func() {
			if !env.IntGreaterThan("GREATER_THAN_INT", 0, 99) {
				t.Error("IntGreaterThan(100 > 99) returned false, want true")
			}
			if env.IntGreaterThan("GREATER_THAN_INT", 0, 100) {
				t.Error("IntGreaterThan(100 > 100) returned true, want false")
			}
		})
	})
}

func TestSystemAndActions(t *testing.T) {
	t.Run("Set and Unset", func(t *testing.T) {
		key := "TEST_SET_UNSET"

		// Test Set
		err := env.Set(key, "value")
		if err != nil {
			t.Fatalf("env.Set returned an error: %v", err)
		}
		val, ok := os.LookupEnv(key)
		if !ok || val != "value" {
			t.Fatal("env.Set failed to set the environment variable correctly")
		}

		// Test Unset
		err = env.Unset(key)
		if err != nil {
			t.Fatalf("env.Unset returned an error: %v", err)
		}
		_, ok = os.LookupEnv(key)
		if ok {
			t.Fatal("env.Unset failed to unset the environment variable")
		}
	})

	t.Run("WasSet and WasUnset", func(t *testing.T) {
		key := "TEST_WAS_SET_UNSET"
		// Ensure the key is not set before starting
		_ = os.Unsetenv(key)

		// Test WasSet
		if !env.WasSet(key, "value") {
			t.Fatal("env.WasSet returned false on successful set")
		}
		val, ok := os.LookupEnv(key)
		if !ok || val != "value" {
			t.Fatal("env.WasSet failed to set the environment variable correctly")
		}

		// Test WasUnset
		if !env.WasUnset(key) {
			t.Fatal("env.WasUnset returned false on successful unset")
		}
		_, ok = os.LookupEnv(key)
		if ok {
			t.Fatal("env.WasUnset failed to unset the environment variable")
		}
	})

	t.Run("User", func(t *testing.T) {
		// This is a basic sanity check. It's difficult to test the error
		// condition without running tests as a user that doesn't exist.
		u := env.User()
		currentUser, err := user.Current()
		if err != nil {
			// If we can't get the current user, we can't reliably test.
			// The error path in env.User() returns a default.
			t.Logf("Could not get current user to compare against: %v. Checking for default.", err)
			if u.Username != "unknown" {
				t.Errorf("env.User() did not return the default user on error, got: %s", u.Username)
			}
			return
		}

		if u == nil {
			t.Fatal("User() returned nil")
		}
		if u.Uid != currentUser.Uid {
			t.Errorf("User() returned incorrect user. Got UID %s, want %s", u.Uid, currentUser.Uid)
		}
		if u.Username != currentUser.Username {
			t.Errorf("User() returned incorrect user. Got username %s, want %s", u.Username, currentUser.Username)
		}
	})
}

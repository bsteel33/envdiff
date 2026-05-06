// Package walker discovers .env files within a directory tree.
//
// It supports configurable glob patterns so callers can match
// standard naming conventions (.env, .env.local, .env.production, etc.)
// as well as custom schemes.
//
// Example usage:
//
//	paths, err := walker.Walk(".", walker.Options{
//		Patterns: []string{".env", ".env.*"},
//		MaxDepth: 2,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range paths {
//		fmt.Println(p)
//	}
package walker

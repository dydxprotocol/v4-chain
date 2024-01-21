// Originally forked from https://github.com/cosmos/iavl/blob/v1.0.0/cmd/iaviewer/main.go
package main

import (
	"bytes"
	"cosmossdk.io/log"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	dbm "github.com/cosmos/cosmos-db"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cosmos/iavl"
	"github.com/spf13/cobra"
)

// TODO: make this configurable?
const (
	DefaultCacheSize int = 10000
)

var (
	prefixRe *regexp.Regexp = regexp.MustCompile("(s/k:[[:alpha:]]+/).*")
)

var rootCmd = &cobra.Command{
	Use:          "iaviewer",
	SilenceUsage: true,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "iaviewer got error: '%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(DataAllCmd())
	rootCmd.AddCommand(DataCmd())
	rootCmd.AddCommand(ShapeCmd())
	rootCmd.AddCommand(VersionsCmd())
}

func AddPathFlag(cmd *cobra.Command) {
	cmd.Flags().String("path", "", "path to leveldb dir")
	if err := cmd.MarkFlagRequired("path"); err != nil {
		panic(err)
	}
}

func AddPrefixFlag(cmd *cobra.Command) {
	cmd.Flags().String("prefix", "", "prefix of the db e.g. \"s/k:gov/\"")
	if err := cmd.MarkFlagRequired("prefix"); err != nil {
		panic(err)
	}
}

func AddVersionFlag(cmd *cobra.Command) {
	cmd.Flags().Uint64("version", 0, "version of the tree (default: latest version)")
}

func DataAllCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-all",
		Short: "Get data of all iavl trees",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			version, _ := cmd.Flags().GetUint64("version")

			db, err := OpenDB(path)
			if err != nil {
				return err
			}
			prefixes := GetAllPrefixes(db)
			for _, prefix := range prefixes {
				tree, err := ReadTree(db, version, []byte(prefix))
				if err != nil {
					return err
				}
				err = PrintTree(tree, prefix)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	AddPathFlag(cmd)
	AddVersionFlag(cmd)
	return cmd
}

func DataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Get data of a specific iavl tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			prefix, _ := cmd.Flags().GetString("prefix")
			version, _ := cmd.Flags().GetUint64("version")

			db, err := OpenDB(path)
			if err != nil {
				return err
			}
			tree, err := ReadTree(db, version, []byte(prefix))
			if err != nil {
				return err
			}
			err = PrintTree(tree, prefix)
			if err != nil {
				return err
			}
			return nil
		},
	}

	AddPathFlag(cmd)
	AddPrefixFlag(cmd)
	AddVersionFlag(cmd)
	return cmd
}

func ShapeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shape",
		Short: "Get shape of a specific iavl tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			prefix, _ := cmd.Flags().GetString("prefix")
			version, _ := cmd.Flags().GetUint64("version")

			db, err := OpenDB(path)
			if err != nil {
				return err
			}
			tree, err := ReadTree(db, version, []byte(prefix))
			if err != nil {
				return err
			}
			PrintShape(tree)
			return nil
		},
	}

	AddPathFlag(cmd)
	AddPrefixFlag(cmd)
	AddVersionFlag(cmd)
	return cmd
}

func VersionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "List versions of specific iavl tree",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			prefix, _ := cmd.Flags().GetString("prefix")
			version, _ := cmd.Flags().GetUint64("version")

			db, err := OpenDB(path)
			if err != nil {
				return err
			}
			tree, err := ReadTree(db, version, []byte(prefix))
			if err != nil {
				return err
			}
			PrintVersions(tree)
			return nil
		},
	}

	AddPathFlag(cmd)
	AddPrefixFlag(cmd)
	AddVersionFlag(cmd)
	return cmd
}

func OpenDB(dir string) (dbm.DB, error) {
	switch {
	case strings.HasSuffix(dir, ".db"):
		dir = dir[:len(dir)-3]
	case strings.HasSuffix(dir, ".db/"):
		dir = dir[:len(dir)-4]
	default:
		return nil, fmt.Errorf("database directory must end with .db")
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	// TODO: doesn't work on windows!
	cut := strings.LastIndex(dir, "/")
	if cut == -1 {
		return nil, fmt.Errorf("cannot cut paths on %s", dir)
	}
	name := dir[cut+1:]
	db, err := dbm.NewGoLevelDB(name, dir[:cut], nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetAllPrefixes fetches all prefixes for iavl trees in the DB. There's no great way to do this,
// so we just iterate through all keys and look for prefixes that match what we expect from modules.
func GetAllPrefixes(db dbm.DB) []string {
	it, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	defer it.Close()
	// Luckily for us, this iterator goes in alphabetical order
	var prefixes []string
	for ; it.Valid(); it.Next() {
		k := it.Key()
		matches := prefixRe.FindStringSubmatch(string(k))
		if matches != nil {
			if len(prefixes) == 0 || matches[1] != prefixes[len(prefixes)-1] {
				prefixes = append(prefixes, matches[1])
			}
		}
	}
	if err := it.Error(); err != nil {
		panic(err)
	}
	return prefixes
}

func PrintDBStats(db dbm.DB) {
	count := 0
	prefix := map[string]int{}
	itr, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()[:1]
		prefix[string(key)]++
		count++
	}
	if err := itr.Error(); err != nil {
		panic(err)
	}
	fmt.Printf("DB contains %d entries\n", count)
	for k, v := range prefix {
		fmt.Printf("  %s: %d\n", k, v)
	}
}

// ReadTree loads an iavl tree from a db
// If version is 0, load latest, otherwise, load named version
// The prefix represents which iavl tree you want to read. The iaviwer will always set a prefix.
func ReadTree(db dbm.DB, version uint64, prefix []byte) (*iavl.MutableTree, error) {
	if len(prefix) != 0 {
		db = dbm.NewPrefixDB(db, prefix)
	}

	tree := iavl.NewMutableTree(db, DefaultCacheSize, false, log.NewLogger(os.Stdout))
	ver, err := tree.LoadVersion(int64(version))
	fmt.Printf("Tree %s version: %d\n", prefix, ver)
	return tree, err
}

func PrintTree(tree *iavl.MutableTree, prefix string) error {
	fmt.Printf("Tree %s data:\n", prefix)
	tree.Iterate(func(key []byte, value []byte) bool { //nolint:errcheck
		printKey := parseWeaveKey(key)
		digest := sha256.Sum256(value)
		fmt.Printf("  %s\n    %X\n", printKey, digest)
		if treeUnmarshallerRegistry, ok := unmarshallerRegistry[prefix]; ok {
			for keyPrefix, unmarshaller := range treeUnmarshallerRegistry {
				if strings.HasPrefix(string(key), keyPrefix) {
					str := unmarshaller(value)
					fmt.Printf("    %s\n", str)
					break
				}
			}
		}
		return false
	})
	hash := tree.Hash()
	fmt.Printf("Tree %s Hash: %X\n", prefix, hash)
	fmt.Printf("Tree %s Size: %X\n", prefix, tree.Size())
	return nil
}

// parseWeaveKey assumes a separating : where all in front should be ascii,
// and all afterwards may be ascii or binary
func parseWeaveKey(key []byte) string {
	cut := bytes.IndexRune(key, ':')
	if cut == -1 {
		return encodeID(key)
	}
	prefix := key[:cut]
	id := key[cut+1:]
	return fmt.Sprintf("%s:%s", encodeID(prefix), encodeID(id))
}

// casts to a string if it is printable ascii, hex-encodes otherwise
func encodeID(id []byte) string {
	for _, b := range id {
		if b < 0x20 || b >= 0x80 {
			return strings.ToUpper(hex.EncodeToString(id))
		}
	}
	return string(id)
}

func PrintShape(tree *iavl.MutableTree) {
	// shape := tree.RenderShape("  ", nil)
	// TODO: handle this error
	shape, _ := tree.RenderShape("  ", nodeEncoder)
	fmt.Println(strings.Join(shape, "\n"))
}

func nodeEncoder(id []byte, depth int, isLeaf bool) string {
	prefix := fmt.Sprintf("-%d ", depth)
	if isLeaf {
		prefix = fmt.Sprintf("*%d ", depth)
	}
	if len(id) == 0 {
		return fmt.Sprintf("%s<nil>", prefix)
	}
	return fmt.Sprintf("%s%s", prefix, parseWeaveKey(id))
}

func PrintVersions(tree *iavl.MutableTree) {
	versions := tree.AvailableVersions()
	fmt.Println("Available versions:")
	for _, v := range versions {
		fmt.Printf("  %d\n", v)
	}
}

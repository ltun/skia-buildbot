package bt

import (
	"context"

	"cloud.google.com/go/bigtable"
	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/go/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TableConfig maps a table name to a list of column families, describing which
// tables and column InitBigtable should create.
type TableConfig map[string][]string

// InitBigtable takes a list of TableConfigs and creates the given tables and
// column families if they don't exist already.
func InitBigtable(projectID, instanceID string, tableConfigs ...TableConfig) error {
	ctx := context.TODO()

	// Set up admin client, tables, and column families.
	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		return sklog.FmtErrorf("Unable to create admin client: %s", err)
	}

	for _, tConfig := range tableConfigs {
		for tableName, colFamilies := range tConfig {
			// Create the table. Ignore error if it already existed.
			err, code := ErrToCode(adminClient.CreateTable(ctx, tableName))
			if err != nil && code != codes.AlreadyExists {
				return sklog.FmtErrorf("Error creating table %s: %s", err)
			} else {
				sklog.Infof("Created table: %s", tableName)
			}

			// Create the column families. Ignore errors if they already existed.
			for _, colFamName := range colFamilies {
				err, code = ErrToCode(adminClient.CreateColumnFamily(ctx, tableName, colFamName))
				if err != nil && code != codes.AlreadyExists {
					return sklog.FmtErrorf("Error creating column family %s in table %s: %s", colFamName, tableName, err)
				}
			}
		}
	}
	return nil
}

// DeleteTables deletes the tables given in the TableConfig.
func DeleteTables(projectID, instanceID string, tableConfigs ...TableConfig) (err error) {
	ctx := context.TODO()

	// Set up admin client, tables, and column families.
	adminClient, err := bigtable.NewAdminClient(ctx, projectID, instanceID)
	if err != nil {
		return sklog.FmtErrorf("Unable to create admin client: %s", err)
	}
	defer func() {
		if err != nil {
			util.Close(adminClient)
		} else {
			err = adminClient.Close()
		}
	}()

	// Delete all tables if they exist.
	for _, tConfig := range tableConfigs {
		for tableName := range tConfig {
			// Ignore NotFound errors.
			err, code := ErrToCode(adminClient.DeleteTable(ctx, tableName))
			if err != nil && code != codes.NotFound {
				return err
			}
		}
	}
	return nil
}

// ErrToCode returns the error that is passed and a gRPC code extracted from the error.
// If the error did not originate in gRPC the returned code is codes.Unknown.
// See https://godoc.org/google.golang.org/grpc/codes for a list of codes.
func ErrToCode(err error) (error, codes.Code) {
	st, _ := status.FromError(err)
	return err, st.Code()
}

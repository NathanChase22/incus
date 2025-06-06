//go:build linux && cgo && !agent

// Code generated by generate-database from the incus project - DO NOT EDIT.

package cluster

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var imageObjects = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByID = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  WHERE ( images.id = ? )
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProject = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  WHERE ( project = ? )
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProjectAndCached = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  WHERE ( project = ? AND images.cached = ? )
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByProjectAndPublic = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  WHERE ( project = ? AND images.public = ? )
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByFingerprint = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  WHERE ( images.fingerprint = ? )
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByCached = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  WHERE ( images.cached = ? )
  ORDER BY projects.id, images.fingerprint
`)

var imageObjectsByAutoUpdate = RegisterStmt(`
SELECT images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update
  FROM images
  JOIN projects ON images.project_id = projects.id
  WHERE ( images.auto_update = ? )
  ORDER BY projects.id, images.fingerprint
`)

// imageColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the Image entity.
func imageColumns() string {
	return "images.id, projects.name AS project, images.fingerprint, images.type, images.filename, images.size, images.public, images.architecture, images.creation_date, images.expiry_date, images.upload_date, images.cached, images.last_use_date, images.auto_update"
}

// getImages can be used to run handwritten sql.Stmts to return a slice of objects.
func getImages(ctx context.Context, stmt *sql.Stmt, args ...any) ([]Image, error) {
	objects := make([]Image, 0)

	dest := func(scan func(dest ...any) error) error {
		i := Image{}
		err := scan(&i.ID, &i.Project, &i.Fingerprint, &i.Type, &i.Filename, &i.Size, &i.Public, &i.Architecture, &i.CreationDate, &i.ExpiryDate, &i.UploadDate, &i.Cached, &i.LastUseDate, &i.AutoUpdate)
		if err != nil {
			return err
		}

		objects = append(objects, i)

		return nil
	}

	err := selectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"images\" table: %w", err)
	}

	return objects, nil
}

// getImagesRaw can be used to run handwritten query strings to return a slice of objects.
func getImagesRaw(ctx context.Context, db dbtx, sql string, args ...any) ([]Image, error) {
	objects := make([]Image, 0)

	dest := func(scan func(dest ...any) error) error {
		i := Image{}
		err := scan(&i.ID, &i.Project, &i.Fingerprint, &i.Type, &i.Filename, &i.Size, &i.Public, &i.Architecture, &i.CreationDate, &i.ExpiryDate, &i.UploadDate, &i.Cached, &i.LastUseDate, &i.AutoUpdate)
		if err != nil {
			return err
		}

		objects = append(objects, i)

		return nil
	}

	err := scan(ctx, db, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"images\" table: %w", err)
	}

	return objects, nil
}

// GetImages returns all available images.
// generator: image GetMany
func GetImages(ctx context.Context, db dbtx, filters ...ImageFilter) (_ []Image, _err error) {
	defer func() {
		_err = mapErr(_err, "Image")
	}()

	var err error

	// Result slice.
	objects := make([]Image, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = Stmt(db, imageObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Project != nil && filter.Public != nil && filter.ID == nil && filter.Fingerprint == nil && filter.Cached == nil && filter.AutoUpdate == nil {
			args = append(args, []any{filter.Project, filter.Public}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(db, imageObjectsByProjectAndPublic)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"imageObjectsByProjectAndPublic\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(imageObjectsByProjectAndPublic)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Project != nil && filter.Cached != nil && filter.ID == nil && filter.Fingerprint == nil && filter.Public == nil && filter.AutoUpdate == nil {
			args = append(args, []any{filter.Project, filter.Cached}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(db, imageObjectsByProjectAndCached)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"imageObjectsByProjectAndCached\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(imageObjectsByProjectAndCached)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Project != nil && filter.ID == nil && filter.Fingerprint == nil && filter.Public == nil && filter.Cached == nil && filter.AutoUpdate == nil {
			args = append(args, []any{filter.Project}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(db, imageObjectsByProject)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"imageObjectsByProject\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(imageObjectsByProject)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID != nil && filter.Project == nil && filter.Fingerprint == nil && filter.Public == nil && filter.Cached == nil && filter.AutoUpdate == nil {
			args = append(args, []any{filter.ID}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(db, imageObjectsByID)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"imageObjectsByID\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(imageObjectsByID)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Fingerprint != nil && filter.ID == nil && filter.Project == nil && filter.Public == nil && filter.Cached == nil && filter.AutoUpdate == nil {
			args = append(args, []any{filter.Fingerprint}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(db, imageObjectsByFingerprint)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"imageObjectsByFingerprint\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(imageObjectsByFingerprint)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Cached != nil && filter.ID == nil && filter.Project == nil && filter.Fingerprint == nil && filter.Public == nil && filter.AutoUpdate == nil {
			args = append(args, []any{filter.Cached}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(db, imageObjectsByCached)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"imageObjectsByCached\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(imageObjectsByCached)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.AutoUpdate != nil && filter.ID == nil && filter.Project == nil && filter.Fingerprint == nil && filter.Public == nil && filter.Cached == nil {
			args = append(args, []any{filter.AutoUpdate}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(db, imageObjectsByAutoUpdate)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"imageObjectsByAutoUpdate\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(imageObjectsByAutoUpdate)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"imageObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID == nil && filter.Project == nil && filter.Fingerprint == nil && filter.Public == nil && filter.Cached == nil && filter.AutoUpdate == nil {
			return nil, fmt.Errorf("Cannot filter on empty ImageFilter")
		} else {
			return nil, errors.New("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getImages(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getImagesRaw(ctx, db, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"images\" table: %w", err)
	}

	return objects, nil
}

// GetImage returns the image with the given key.
// generator: image GetOne
func GetImage(ctx context.Context, db dbtx, project string, fingerprint string) (_ *Image, _err error) {
	defer func() {
		_err = mapErr(_err, "Image")
	}()

	filter := ImageFilter{}
	filter.Project = &project
	filter.Fingerprint = &fingerprint

	objects, err := GetImages(ctx, db, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"images\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, ErrNotFound
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"images\" entry matches")
	}
}

#!/usr/bin/env bash
cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator" down
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
cd ../..


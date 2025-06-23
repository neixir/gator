#!/usr/bin/env bash
cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator" $1 $2
cd ../..

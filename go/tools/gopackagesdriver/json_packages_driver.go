// Copyright 2021 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"go/types"
)

type JSONPackagesDriver struct {
	registry *PackageRegistry
}

func NewJSONPackagesDriver(jsonFiles []string, prf PathResolverFunc) (*JSONPackagesDriver, error) {
	jpd := &JSONPackagesDriver{
		registry: NewPackageRegistry(),
	}

	for _, f := range jsonFiles {
		if err := WalkFlatPackagesFromJSON(f, func(pkg *FlatPackage) {
			jpd.registry.Add(pkg)
		}); err != nil {
			return nil, fmt.Errorf("unable to walk json: %w", err)
		}
	}

	if err := jpd.registry.ResolvePaths(prf); err != nil {
		return nil, fmt.Errorf("unable to resolve paths: %w", err)
	}

	if err := jpd.registry.ResolveImports(); err != nil {
		return nil, fmt.Errorf("unable to resolve paths: %w", err)
	}

	return jpd, nil
}

func (b *JSONPackagesDriver) Match(pattern ...string) *driverResponse {
	rootPkgs, packages := b.registry.Match(pattern...)

	return &driverResponse{
		NotHandled: false,
		Sizes:      types.SizesFor("gc", "amd64").(*types.StdSizes),
		Roots:      rootPkgs,
		Packages:   packages,
	}
}

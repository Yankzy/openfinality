#!/usr/bin/env bash

set -e

echo "Verifying governance files..."
for file in README.md ROADMAP.md SECURITY.md ARCHITECTURE.md GOVERNANCE.md CONTRIBUTING.md CODE_OF_CONDUCT.md SUPPORT.md MAINTAINERS.md CHANGELOG.md VERSIONING.md RELEASE.md; do
    if [ ! -f "$file" ]; then
        echo "Missing required governance file: $file"
        exit 1
    fi
done

echo "Checking for tracked private keys..."
if git grep -q -e "-----BEGIN PRIVATE KEY-----"; then
    echo "Found tracked private key!"
    exit 1
fi

echo "Checking for placeholders (TODO/FIXME)..."
# Ignore this file itself and the demo scripts if they contain the words
if git grep -E "(TODO|FIXME|panic\(\"not implemented\"\)|return nil // placeholder)" -- \
    ':!tools/verify-repo.sh' ':!docs/*' ':!.github/*'; then
    echo "Found placeholders in source code!"
    exit 1
fi

echo "Checking OpenAPI validation..."
if [ ! -f "protocol/openapi/openapi.yaml" ]; then
    echo "OpenAPI spec missing"
    exit 1
fi

echo "Verification complete. Afro-Rail repo is clean."
exit 0

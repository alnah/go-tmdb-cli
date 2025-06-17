#!/bin/bash

# Script pour ajouter t.Parallel() à tous les tests

add_parallel() {
    local file=$1
    echo "Processing $file..."

    # Utiliser sed pour ajouter t.Parallel() après chaque déclaration de fonction Test
    # Pour macOS, utiliser -i '' pour édition en place
    sed -i '' '/^func Test[A-Za-z_]*([[:space:]]*t[[:space:]]*\*testing\.T)[[:space:]]*{$/a\
\	t.Parallel()
' "$file"
}

# Ajouter aux nouveaux fichiers de test
for file in tests/unit/client_movies_test.go tests/unit/client_tv_test.go tests/unit/client_search_test.go; do
    if [ -f "$file" ]; then
        add_parallel "$file"
    fi
done

echo "t.Parallel() added to test files!"

---
Title: Scripts
StyleSheets:
  - /static/styles/utils/code.css
Scripts:
  - /static/scripts/utils/copy-code.js
---

# Scripts

## Docker

### dc

```bash
#!/bin/bash

docker-compose $@
```

## Git

### ga

`git add` but with the ability to add all files in current dir by default

```bash
#!/bin/bash
if [ "$#" -eq  "0" ]
  then
    git add .
else
    git add $@
fi
```

### gs

```bash
#!/bin/bash

git status
```

## ls

### l

Lists the files in current dir with file sizes in human readable format

```bash
#!/bin/bash

ls -lh
```

### la

Lists all files (including hidden) in current dir

```bash
#!/bin/bash

ls -la
```

### lg

Lists all files (including hidden) in current dir sorted by size measured as number of files

```bash
#!/bin/bash

ls -lah | sort -k1 -r
```

### ll

Lists all files (including hidden) in current dir with file sizes in human readable format

```bash
#!/bin/bash

ls -lah
```

### lm

Lists files in current dir with comma separated values

```bash
#!/bin/bash

ls -m
```

### lr

Lists files in current dir and all subdirectories recursively

```bash
#!/bin/bash

ls -R
```

### rmrf

Removes files and directories recursively

```bash
#!/bin/bash

rm -rf $@
```

# Mugit

Mugit is a **Git-like version control system written from scratch in Go**.

The project is being developed to understand how version control systems work internally rather than treating Git as a black box. Mugit progressively implements the core concepts behind Git, including repository initialization, content-addressable storage, object hashing, blobs, branches, commits, and eventually merging and remote repositories.

The project focuses on understanding the **data structures, storage model, algorithms, and design decisions** behind a distributed version control system.

---

## Why Mugit?

Git is used every day by software engineers, but many of its internal mechanisms are hidden behind simple commands such as:

```bash
git add
git commit
git branch
git merge
git push
```

Mugit recreates these concepts from the ground up to understand what actually happens behind these commands.

The main goals are to:

* Understand how version control systems store data
* Learn content-addressable storage
* Understand Git's object model
* Understand commits and commit history
* Understand branches and `HEAD`
* Learn how merging and conflicts work
* Practice systems programming with Go
* Build a non-trivial engineering project from first principles

---

## Current Features

### Repository Initialization

Mugit can initialize a repository in the current directory:

```bash
mugit init
```

This creates the internal `.mugit` directory:

```text
.mugit/
├── objects/
├── refs/
│   └── heads/
├── index
└── HEAD
```

The `.mugit` directory contains the internal metadata and objects required by Mugit.

---

### Repository Discovery

Mugit can locate the repository root by walking upward through the filesystem.

For example:

```text
project/
├── .mugit/
└── src/
    └── services/
        └── auth/
```

Running a Mugit command from `auth/` can search upward until it finds:

```text
.mugit/
```

This allows commands to work from nested directories inside the repository.

---

### Repository Status

Mugit currently supports a basic status command:

```bash
mugit status
```

Example:

```text
On branch main
No commits yet
```

Once commits exist, Mugit will report the current commit associated with the branch.

---

### Content-Addressable Object Storage

Mugit uses SHA-1 hashes to identify objects.

A blob is represented using a Git-like format:

```text
blob <size>\0<content>
```

For example:

```text
blob 5\0Hello
```

The SHA-1 hash of this representation becomes the object's identifier.

This means:

```text
Content
   ↓
Blob representation
   ↓
SHA-1
   ↓
Object ID
   ↓
Stored object
```

The same content produces the same object ID, allowing identical objects to be reused instead of stored multiple times.

---

### `hash-object`

Mugit can create and store a blob object from a file:

```bash
mugit hash-object file.txt
```

The command produces an object ID such as:

```text
5d41402abc4b2a76b9719d911017c592
```

The object is stored inside:

```text
.mugit/objects/
```

Existing objects are not unnecessarily rewritten.

---

### `cat-file`

Mugit can retrieve the content of a stored object:

```bash
mugit cat-file <object-id>
```

For example:

```text
Hello
```

Mugit reads the stored object, separates its header from its content, and returns the actual blob content.

---

## Architecture

The project currently follows a simple architecture:

```text
CLI
 │
 ├── init
 ├── status
 ├── hash-object
 └── cat-file
       │
       ▼
Repository Layer
       │
       ▼
.mugit/
       │
       ├── objects/
       ├── refs/
       │   └── heads/
       ├── index
       └── HEAD
```

The implementation is intentionally kept simple while the underlying concepts are being developed.

As the project grows, the codebase will be refactored into clearer modules and layers.

---

## Planned Features

Mugit is being developed incrementally.

### Core Version Control

* [x] Repository initialization
* [x] Repository root discovery
* [x] `status`
* [x] Blob objects
* [x] SHA-1 object IDs
* [x] Object storage
* [x] `hash-object`
* [x] `cat-file`
* [ ] Object validation
* [ ] Index / staging area
* [ ] `add`
* [ ] `reset`
* [ ] Tree objects
* [ ] `commit`
* [ ] `log`
* [ ] `show`

### Branching

* [ ] Branch creation
* [ ] Branch listing
* [ ] Branch switching
* [ ] `HEAD` management
* [ ] Detached `HEAD`

### History & Changes

* [ ] `diff`
* [ ] Staged diff
* [ ] Commit diff
* [ ] Three-way merge
* [ ] Merge conflicts
* [ ] `merge`

### Advanced Version Control

* [ ] Tags
* [ ] `revert`
* [ ] `reset`
* [ ] Object integrity checking
* [ ] Garbage collection
* [ ] Compression
* [ ] Binary object format
* [ ] Remote repositories
* [ ] `fetch`
* [ ] `push`
* [ ] `pull`

### Future Exploration

Potential future directions include:

* Commit graph visualization
* Web interface
* Repository statistics
* Change-impact analysis
* Semantic merge assistance
* Dependency-aware history
* Distributed collaboration
* Git compatibility

---

## Technology

* **Language:** Go
* **Version Control:** Git
* **Repository Hosting:** GitHub

The project intentionally uses Go to practice:

* File system operations
* Binary data handling
* Hashing
* CLI development
* Data structures
* Error handling
* Systems programming
* Repository and storage design

---

## Learning Approach

Mugit is being built incrementally rather than copied from an existing Git implementation.

Each feature follows a development cycle:

```text
Understand the concept
        ↓
Design the behavior
        ↓
Implement
        ↓
Test
        ↓
Inspect failures
        ↓
Refactor
        ↓
Move to the next subsystem
```

The goal is not to reproduce every feature of Git, but to develop a deep understanding of the engineering concepts that make a version control system possible.

---

## Project Status

🚧 **Active Development**

Mugit is currently in the early implementation stage.

The repository and blob object foundations are implemented, and development is progressing toward the staging/index system, tree objects, and commits.

---

## License

This project is currently intended as a learning and portfolio project.

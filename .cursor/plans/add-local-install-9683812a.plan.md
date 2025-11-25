<!-- 9683812a-c337-440e-8da9-e1d530e25fce f1707f1e-ff84-4374-b3fd-8dfa55f80dcf -->
# Add Local Install Command to Makefile

## Overview

Add a new Makefile target that builds the binary locally and installs it to `/usr/local/bin/` for use with Cursor and other MCP clients.

## Changes Required

### 1. Add `install-local` target to Makefile

Add a new target after the existing `install` target:

```makefile
## install-local: Build and install binary to /usr/local/bin (requires sudo)
install-local: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Binary installed to /usr/local/bin/$(BINARY_NAME)"
	@echo "You can now use '$(BINARY_NAME)' command globally"
```

This target will:

- First build the binary (via `build` dependency)
- Copy the binary to `/usr/local/bin/` (requires sudo)
- Make it executable
- Provide confirmation messages

### 2. Update the help text

The new target will automatically appear in `make help` output due to the `##` comment format.

## Usage

After implementation, users can run:

```bash
make install-local
```

This will build and install the binary system-wide for use with Cursor's MCP configuration.

### To-dos

- [ ] Add install-local target to Makefile after the install target
- [ ] Verify the new target appears in make help output
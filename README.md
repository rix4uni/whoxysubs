## whoxysubs

Scrape whoxy subdomains without api key.

## Installation
```
go install github.com/rix4uni/whoxysubs@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/whoxysubs/releases/download/v0.0.1/whoxysubs-linux-amd64-0.0.1.tgz
tar -xvzf whoxysubs-linux-amd64-0.0.1.tgz
rm -rf whoxysubs-linux-amd64-0.0.1.tgz
mv whoxysubs ~/go/bin/whoxysubs
```
Or download [binary release](https://github.com/rix4uni/whoxysubs/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/whoxysubs.git
cd whoxysubs; go install
```

## Usage
```
Usage of whoxysubs:
  -s, --search string   Search type: company, email, keyword, or name
      --silent          silent mode.
      --version         Print the version of the tool and exit.
```

## Usage Examples

1. **Search by Company:**
   ```bash
   echo "Dell Inc." | go run main.go -s company
   ```

2. **Search by Email:**
   ```bash
   echo "dns-admin@google.com" | go run main.go -s email
   ```

3. **Search by Keyword:**
   ```bash
   echo "dell.com" | go run main.go -s keyword
   ```

4. **Search by Name:**
   ```bash
   echo "elon musk" | go run main.go -s name
   ```
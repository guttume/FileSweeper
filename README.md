# File Sweeper

The File Sweeper is a Go program designed to automate the management of files and directories based on specified criteria such as age and action.

## Features

- Monitor multiple locations for files and directories.
- Define cleanup actions including deletion and moving of files.
- Specify criteria for file cleanup based on the number of days since modification.
- Log cleanup actions for auditing and troubleshooting purposes.

## Usage

### 1. Installation

Clone the repository:

```bash
git clone <repository_url>
cd file-cleanup-utility
```

### 2. Configuration
Create a configuration file (config.json by default) to define the locations to monitor and cleanup actions. Here's an example configuration:

```json
{
  "locations": [
    {
      "path": "/path/to/source/directory",
      "days": 30,
      "action": "delete",
      "target": "/path/to/target/directory"
    }
  ],
  "log_file": "logs/app.log"
}
```

- locations: An array of objects representing the directories to monitor.
    - path: The path to the directory to monitor.
    - days: The number of days since the last modification after which files should be considered for cleanup.
    - action: The cleanup action to perform. Supported actions are "delete" and "move".
    - target: (Optional) The target directory for files to be moved. Required only if action is set to "move".
- log_file: The path to the log file to store cleanup actions.

### 3. Execution
Run the program with the following command:

```bash
go run main.go -c config.json
```

Replace config.json with the path to your configuration file if it's located elsewhere.

## License
This project is licensed under the MIT License.

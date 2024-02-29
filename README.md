# Zoho Sign File Downloader

## Introduction

This script, developed in Go and leveraging multi-threading capabilities, is designed to efficiently connect with the Zoho Sign API and download all files uploaded by customers. It's intended for businesses or individuals looking to automate the retrieval of documents submitted via Zoho Sign, ensuring a streamlined process for managing and storing signed documents securely.

## Features

- **Multi-threading Support**: Utilizes Go's concurrency model to handle multiple downloads simultaneously, improving efficiency.
- **Zoho Sign API Integration**: Securely connects with Zoho Sign API to access and download customer-uploaded files.
- **Automatic Retry Mechanism**: Implements a retry logic for failed downloads, ensuring robustness in network variability conditions.
- **Customizable Download Directory**: Allows users to specify a download directory for organizing the retrieved files.

## Prerequisites

Before you begin, ensure you have the following:

- Go installed on your machine (version 1.15 or later recommended).
- A Zoho Sign API key, obtainable through your Zoho Sign account for authenticating API requests.
- Basic knowledge of Go programming and understanding of concurrent programming concepts.

## Usage
To extract the zip file and merge the PDF, please use extractor.py. This utility is designed to complement the downloading process by providing an easy way to consolidate downloaded documents into a single PDF file for easier handling and storage.
## Support and Contributions
For support requests, feature suggestions, or contributions, please open an issue in the GitHub repository or submit a pull request. Your feedback and contributions are welcome to improve the script's functionality and reliability.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

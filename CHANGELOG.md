# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.1] - 2025-12-05

### 🚀 Performance Improvements
- **Migrated Core Detectors to Custom Parsers:**
  - Removed `regexp` dependency for **Email**, **Credit Card**, **IPv4**, **Phone**, **CPF**, and **CNPJ**.
  - Replaced `regexp` in `Restore()` method with a single-pass string builder loop.
  - **10x speedup** on average compared to previous regex implementations.
  - Drastically reduced GC pressure (allocations are now limited to result slice creation).

### ✨ Features
- **Enhanced IP Detector:** Added strict boundary checks to prevent false positives in version numbers (e.g., `1.2.3`), negative numbers, and leading zeros.
- **Enhanced Email Detector:** Improved "left/right expansion" logic to correctly identify emails embedded in complex strings.
- **Enhanced Phone Detector:** Added strict checks for E.164 format and runaway loop protection.

### 🐛 Bug Fixes
- Fixed false positives in UUID detection where long hex strings were captured as UUIDs.
- Fixed boundary issues in CPF/CNPJ parsers.

---

## [v1.0.0] - 2025-12-04
- Initial Release.
# ChronoGrid (University Timetable Studio)

> **Status:** Work in Progress (Core visualization, filtering, and styled Excel export engine implemented).

A modern web application for university course scheduling, room allocation tracking, and formatted timetable exporting—built with **Angular (Signals & Standalone Components)**, **Material Design 3**, and **SheetJS (`xlsx-js-style`)**.

## Key Features
- **Multi-Perspective Views:** Instant switching between **Academic Levels** (with group filtering), **Classrooms**, and **Professors**.
- **High-Performance Grid Engine:** Decouples full matrix construction (cached in memory) from zero-latency numeric filtering.
- **Custom Session Intervals:** Configurable up to 10 time slots with `localStorage` persistence.
- **Theme-Aware Excel Export (`.xlsx`):** Dynamically extracts active Material 3 CSS tokens at runtime to generate styled, print-ready spreadsheets.
- **Landscape Print Layout:** Dedicated print stylesheet with automatic text wrapping and zero clipping.

## Roadmap
- [ ] Interactive "What-If" Sandbox for manual drag/click session rescheduling.
- [ ] Real-time free-rooms availability matrix and conflict detection.
- [ ] Automated constraint-based timetable solver integration.
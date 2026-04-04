# Stack Explorer Frontend — Sigil Project

*Build a polished, modern frontend for the Stack Explorer research platform using Sigil.*

## Goal

Create a full-featured web UI for Stack Explorer that provides CRUD over all entities, interactive reports, and a dashboard that showcases what Sigil can do. This is a flagship Sigil project — it should look and feel like a premium analytics tool, not a basic admin panel.

**This is a push-the-boundaries project.** Use every component type that fits. Make the layouts sophisticated. The UX should feel like a modern SaaS analytics dashboard — think Linear meets Datadog.

## Product Vision

Read `/Users/chrispian/Projects-apps/stack-explorer/docs/product-vision.md` for full context on where the product is headed. Key points:

- Three consumer tracks: Engineer, Product, Leadership
- Multi-lens scoring — same data, different perspectives
- Time-series trends, gap analysis, pattern discovery
- Eventually: CVE tracking, job market data, ecosystem graphs

## Stack Explorer Data Model

The SQLite database at `stack-explorer/data/stack-explorer.db` has these tables:

| Table | Purpose | Key Fields |
|-------|---------|------------|
| `repos` | Repo catalog (112 repos) | id, name, url, description, stack, category, is_own, local_path |
| `tags` / `repo_tags` | Many-to-many tagging | tag.id, tag.name |
| `snapshots` | Time-series metrics | repo_id, captured_at, stars, forks, loc, files, contributors, commits_30d, complexity_avg |
| `review_dimensions` | 18 scoring axes | id, name, category, weight, description |
| `lenses` | 9 scoring perspectives | id, name, description |
| `lens_dimensions` | Lens-dimension weights | lens_id, dimension_id, weight |
| `scorecards` | Per-repo score containers | repo_id, lens_id, overall, scored_at |
| `dimension_scores` | Individual scores | scorecard_id, dimension_id, score, evidence |
| `architecture_patterns` | 26 patterns/anti-patterns | id, name, type, category, description |
| `repo_patterns` | Pattern-repo links | repo_id, pattern_id, quality, notes |
| `findings` | Research observations | repo_id, title, category, severity, status, description |
| `code_references` | File:line pointers | repo_id, file_path, line_start, description, pattern_id, finding_id |
| `comparison_sets` / `comparison_set_repos` | Gap analysis groups | set_id, repo_id, role (subject/reference) |
| `report_configs` / `report_config_filters` | Report definitions | config_id, lens_id, filter_type, filter_value |

## API Layer

Stack Explorer currently has a CLI only. The frontend will need a REST API. Options:
1. **Add an HTTP server to the Go CLI** (preferred — `internal/api/` with Chi router, same as cerberus/conduit)
2. **Generate handler stubs from Sigil datasources** and implement them

The API should follow the Sigil datasource endpoint contract:
- `GET /api/<resource>?search=&filter[field]=&sort=&direction=&page=&pageSize=`
- `POST /api/<resource>` → 201
- `GET /api/<resource>/{id}` → 200
- `PUT /api/<resource>/{id}` → 200
- `DELETE /api/<resource>/{id}` → 204

## Pages to Build

### 1. Dashboard (`dashboard`)
The landing page. At-a-glance overview of the research landscape.

- **Stat cards row**: Total repos, scored repos, dimensions, patterns, findings
- **Chart**: Score distribution by category (bar chart)
- **Chart**: Own projects vs reference averages (radar/bar by dimension)
- **Top movers**: Repos with biggest score changes (when we have multiple snapshots)
- **Recent findings**: Latest findings with severity badges
- **Quick actions**: Add repo, run scan, generate report

### 2. Repo Catalog (`repos`)
Full CRUD for the repo catalog. The primary data management view.

- **Data table** with columns: name (linked), category, stack, own (badge), score, LoC, tags
- **Search bar** with debounce
- **Filter chips**: by category, stack, is_own, tag
- **Toolbar**: Add Repo (modal), Import YAML, Export YAML, Bulk Actions
- **Row click** → Repo Detail page
- **Add/Edit modals** with all repo fields + tag management
- **Bulk delete** with confirmation

### 3. Repo Detail (`repo-detail`)
Deep dive into a single repo. Overlay: sheet or full page.

- **Header**: Name, URL (linked), category badge, stack badge, own indicator
- **Tabs**:
  - **Scores**: Dimension scores table with evidence, lens selector dropdown to recalculate overall
  - **Snapshots**: Time-series table (or chart) of LoC, stars, commits over time
  - **Patterns**: Linked patterns with quality badges
  - **Findings**: Findings for this repo
  - **Tags**: Tag chips with add/remove
- **Actions**: Edit, Score, Scan, Delete

### 4. Scorecards (`scorecards`)
View and compare scores across repos.

- **Lens selector** dropdown at top — changes all scores on the page
- **Data table**: Repo, Overall, then one column per dimension in the selected lens
- **Heatmap coloring**: Green (8+), yellow (5-7), red (<5), gray (unscored)
- **Compare mode**: Select 2-3 repos → side-by-side radar chart
- **Generate scorecard**: Button to export markdown

### 5. Reports (`reports`)
Manage and generate reports from report configs.

- **Data table**: Report name, lens, repo count, audience, actions
- **Create Report modal**: Name, lens selector, filter builder (add multiple type=value filters)
- **Generate button** → Downloads or displays markdown
- **Report detail**: Shows the generated leaderboard with dimension breakdown table

### 6. Gap Analysis (`gap-analysis`)
Visual gap analysis tool.

- **Comparison set selector** dropdown
- **Split view**: Subjects on left, reference average on right
- **Bar chart**: Delta per dimension (red=gap, green=advantage)
- **Table**: Dimension, subject score, ref avg, delta, status badge
- **Create comparison set** modal with repo multi-select (subjects vs references)
- **Findings**: Auto-generated gap findings linked below

### 7. Patterns (`patterns`)
Architecture pattern catalog.

- **Data table**: Pattern name, type (pattern/anti-pattern badge), category, linked repos count
- **Add Pattern modal**: Name, type, category, description
- **Pattern detail**: Description + linked repos with quality badges + code references
- **Link Pattern modal**: Select repo, quality rating, notes

### 8. Findings (`findings`)
Research findings tracker.

- **Data table**: Title, repo, category badge, severity badge, status badge
- **Filter chips**: category (gap/strength/opportunity/risk), severity, status
- **Add Finding modal**: Title, repo selector, category, severity, description
- **Finding detail**: Full description + linked code references + linked patterns

### 9. Dimensions & Lenses (`dimensions`)
Configuration management for the scoring system.

- **Tabs**: Dimensions | Lenses
- **Dimensions tab**: Data table with name, category, weight, description. Add Dimension modal.
- **Lenses tab**: Data table with name, dimension count, description. Lens Detail shows dimensions with weights. Create Lens modal with dimension multi-select and weight inputs.

### 10. Settings (`settings`)
Application configuration.

- **Database stats** (table row counts)
- **Import/Export** buttons (repos YAML, full DB backup)
- **Scan configuration** (blueprint paths)
- **About** section with version

## DataSources to Define

Create Sigil datasource manifests for each entity:

| Alias | Endpoint Base | Capabilities |
|-------|--------------|-------------|
| `Repo` | `/api/repos` | search, filter, sort, paginate, create, read, update, delete |
| `Tag` | `/api/tags` | search, paginate, create, delete |
| `Snapshot` | `/api/snapshots` | filter, sort, paginate, create |
| `Dimension` | `/api/dimensions` | search, paginate, create |
| `Lens` | `/api/lenses` | search, paginate, create, read, delete |
| `Scorecard` | `/api/scorecards` | filter, sort, paginate, read |
| `DimensionScore` | `/api/scores` | filter, create, update |
| `Pattern` | `/api/patterns` | search, filter, sort, paginate, create, read, update, delete |
| `Finding` | `/api/findings` | search, filter, sort, paginate, create, read, update |
| `ComparisonSet` | `/api/comparison-sets` | search, paginate, create, read, delete |
| `ReportConfig` | `/api/reports` | search, paginate, create, read, delete |

## Theme

Use a dark theme inspired by the Hadron desktop app's design language:
- Background: zinc-950 (#09090b)
- Surface: zinc-900 (#18181b) 
- Borders: zinc-800 (#27272a)
- Accent: blue-500 for primary actions and links
- Success/green for high scores, amber for medium, red for low/critical
- Use the blue glass hover effect from Hadron's BlueprintsPage pattern

## Design Principles

1. **Information density over whitespace.** This is an analytics tool — pack it with useful data.
2. **Keyboard-first.** Cmd+K for command palette, shortcuts for common actions.
3. **Lens context everywhere.** The current lens should be visible and switchable globally.
4. **Scores are visual.** Use color coding, progress bars, or heatmaps — never just numbers in a table.
5. **Drill-down.** Everything clickable leads deeper. Repo → scores → dimension → evidence.
6. **Responsive.** Works on laptop screens (1440px primary target).

## Implementation Order

1. **DataSources first** — Define all 11 datasource manifests
2. **Theme** — Create a dark theme matching the design spec
3. **Dashboard** — Landing page with stat cards and charts
4. **Repos** — Full CRUD catalog (this exercises the most Sigil features)
5. **Scorecards** — The signature view with lens switching and heatmaps
6. **Gap Analysis** — Charts and comparison views
7. **Reports** — Config management and generation
8. **Patterns & Findings** — CRUD with linking
9. **Dimensions & Lenses** — Configuration pages
10. **Settings** — Utility page

## API Implementation Note

After generating the frontend, the Go API layer needs to be built in `stack-explorer/internal/api/`. Follow the cerberus pattern:
- Chi router
- JSON responses matching the Sigil datasource contract
- Reuse the existing `internal/store/sqlite/` layer
- No auth needed initially (local tool)

## Reference Files

- Stack Explorer CLAUDE.md: `/Users/chrispian/Projects-apps/stack-explorer/CLAUDE.md`
- Product vision: `/Users/chrispian/Projects-apps/stack-explorer/docs/product-vision.md`
- Stack Explorer DB: `/Users/chrispian/Projects-apps/stack-explorer/data/stack-explorer.db`
- Sigil component schemas: `/Users/chrispian/Projects-apps/sigil/internal/components/builtin/`
- Sigil datasource spec: `/Users/chrispian/Projects-apps/sigil/docs/06_datasource-model.md`
- Sigil config spec: `/Users/chrispian/Projects-apps/sigil/docs/02_config-spec.md`
- Hadron design language: `/Users/chrispian/Projects-apps/hadron/.agentrc/boot-prompt.md` (design section)
- Example page: `/Users/chrispian/Projects-apps/sigil/examples/sprint-dashboard.yaml`

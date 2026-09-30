# Sigil UI Compiler Direction

**Status:** Directional architecture draft

**Date:** 2026-08-22
**Scope:** Desired product boundaries and architecture; not an implementation plan

## Purpose

This document describes the intended place of Sigil in the Hollis Labs
portfolio. It evaluates Sigil's definition language, compiler pipeline,
component registry, data and action contracts, themes, renderers, generated
artifacts, authoring interfaces, preview server, MCP surface, extensions,
security, and portfolio composition.

It deliberately does not define phases, estimates, task breakdowns, migration
steps, release sequencing, or compatibility work. Existing documentation and
code remain the current implementation truth until explicitly superseded.
This document defines a target direction for a later architecture and planning
session.

## Portfolio axiom

> Hollis tools own execution and operational state, but not the business
> definitions or business data they operate on.

For Sigil, this means:

- A project or UI publisher owns its page, application, module, theme,
  datasource, action, and custom-component definitions.
- A design-system publisher owns its component implementations, public APIs,
  design tokens, accessibility behavior, and release lifecycle.
- A backend or schema publisher owns the API and data contracts to which a UI
  binds.
- A consuming repository owns the generated artifacts after Sigil applies
  them under an explicit generation authority mode.
- The consuming application owns runtime rendering, authorization, data,
  effects, navigation, and user state.
- Credential authorities own secrets and their lifecycle.
- Sigil owns compiler semantics and compilation facts: source discovery,
  parsing, normalization, reference resolution, validation, intermediate
  representation, renderer compatibility, deterministic artifact production,
  safe output reconciliation, diagnostics, provenance, and verification
  reports.

Loading a UI definition does not make Sigil its author. Generating a handler
stub does not make Sigil the backend. Generating an HTTP action does not make
Sigil the effect authority. Copying a custom component does not transfer its
design-system ownership to Sigil.

The useful shorthand is:

> Sigil compiles publisher-owned UI intent into deterministic,
> framework-native source artifacts without becoming the application's
> runtime or business authority.

## Product definition

Sigil is a system-agnostic, build-time UI compiler. It accepts a versioned set
of declarative UI definitions, resolves them against registered component and
data contracts, compiles them into a canonical UI intermediate representation,
and materializes framework-native source artifacts through a selected renderer.

Sigil:

1. Discovers and parses a declared UI source bundle.
2. Resolves versioned component, theme, data, route, and extension references.
3. Validates portable semantics and target-specific capabilities.
4. Compiles source definitions into one canonical, framework-neutral UI IR.
5. Invokes a pinned renderer against that IR.
6. Produces a deterministic artifact set with ownership and provenance.
7. Safely reconciles managed artifacts or transfers scaffold artifacts to the
   consumer according to an explicit authority mode.
8. Optionally formats, builds, type-checks, and verifies the result through
   declared toolchain adapters.
9. Exposes the same compiler and authoring semantics through CLI, Go API, MCP,
   and local preview adapters.

Sigil is not a runtime renderer, low-code application platform, visual editor,
backend framework, API gateway, database mapper, authorization system,
workflow engine, agent runtime, infrastructure control plane, secret manager,
or general-purpose templating language.

Generated applications do not require Sigil at runtime.

## Core product boundary: compiler, not runtime

The compiler boundary is a durable product decision.

Sigil owns build-time operations:

- parsing and validating UI definitions
- resolving component and data contracts
- compiling portable UI semantics
- selecting renderer capabilities
- generating and reconciling source artifacts
- creating static previews and catalogs as build artifacts
- reporting drift, compatibility, provenance, and verification

The generated application owns runtime operations:

- loading and mutating user or business data
- authenticating users and enforcing authorization
- performing HTTP, navigation, clipboard, file, and other effects
- maintaining browser, desktop, or server state
- rendering compiled components to end users
- reporting application telemetry and failures

`sigil serve` is a local development projection over source definitions. It
does not turn Sigil into the production UI runtime. The preview renderer is a
compiler target with deliberately limited semantics, not an alternate
application platform.

Tangent remains a runtime interaction host. Sigil may scaffold or generate
Tangent renderer source at build time, but Tangent never invokes Sigil to
render an envelope at runtime.

## Architecture sketch

```text
UI publisher repositories                     Contract publishers
pages · app/modules · themes                   design system · API schema
datasources · custom components                component packs · renderer packs
                    |                                      |
                    +--------------+-----------------------+
                                   |
                            SourceBundle
                                   v
                      +--------------------------+
                      |          Sigil           |
                      |                          |
 CLI / Go API ------> | source loader            |
 MCP ---------------> | resolver + validator     |
 preview server ----> | canonical UI compiler    |
                      | capability evaluator     |
                      | renderer host            |
                      | artifact reconciler      |
                      | verification adapters    |
                      +------------+-------------+
                                   |
                      CompilationReceipt
                      ArtifactSet + ownership manifest
                                   |
                +------------------+------------------+
                |                                     |
         managed generation                    scaffold generation
         reproducible projection               explicit ownership transfer
                |                                     |
                +------------------+------------------+
                                   v
                         consuming application
                    build · test · release · runtime

Catalog and preview are output projections. Studio, if present, is a sibling
authoring application over Sigil's public compiler contract.
```

The source definition, compiled IR, renderer registration, staged artifact
set, applied files, verified build, released application, and running UI are
different objects with different authorities.

## Responsibility boundaries

| Concern | Authoritative owner | Sigil responsibility |
|---|---|---|
| UI purpose and business behavior | Consuming application or UI publisher | Preserve declared intent; do not invent domain semantics |
| Page and application definitions | UI publisher | Parse, validate, resolve, compile, and report exact revisions |
| Component semantic contract | Component or design-system publisher | Register, version, resolve, and validate use |
| Component runtime implementation | Design-system package or project | Import, copy, or reference through a declared binding |
| Theme and design tokens | Design-system or project publisher | Resolve, validate, transform, and materialize |
| Backend schema and endpoint behavior | Backend or schema publisher | Bind a pinned UI-side projection and generate compatible clients or stubs |
| Runtime data | Consuming application and backend | Generate access code; never persist or become data authority |
| Runtime authorization | Consuming application and identity system | Generate integration hooks or guards without treating them as enforcement |
| Runtime effects | Consuming application | Generate typed calls or handlers; never execute effects during compilation |
| Portable UI IR | Sigil | Define compiler semantics and preserve source provenance |
| Renderer implementation | Renderer publisher | Validate registration, invoke safely, and record exact version/digest |
| Compilation | Sigil | Produce deterministic artifacts, diagnostics, and receipts |
| Managed generated files | Sigil as projection manager; consumer repository as custodian | Reconcile only declared owned paths and detect divergence |
| Scaffold output | Consumer after explicit transfer | Produce once and record that Sigil no longer manages it |
| Handwritten extension code | Project or package publisher | Reference through explicit extension points; never silently overwrite |
| Build and release acceptance | Consumer build system | Run declared verification adapters and report evidence |
| Authoring UX | CLI, agent client, or sibling Studio | Provide compiler and mutation contracts, not a hidden document authority |
| Infrastructure and local services | Cerberus and target runtime | Expose safe commands, health, and foreground process behavior |
| Secrets | External credential authority | Carry references or environment names; never persist values in definitions or artifacts |

## Authority and custody model

### Source bundle

A compilation consumes an explicitly rooted, immutable source snapshot:

```text
SourceBundle
  project identity and root
  definition-language version
  project and app definitions
  page and shell definitions
  datasource and route bindings
  theme definitions
  component-pack registrations
  custom component and provider source refs
  renderer and target profile
  source file identities and digests
  repository revision and dirty-state disclosure
```

The project repository remains the source of truth. Sigil may hold an in-memory
parse cache or build cache, but it does not create a second mutable catalog of
project definitions.

Source discovery is configuration-driven and rooted. Ambient files outside the
declared source roots do not silently participate in a build.

### Definition revision

Every parsed definition has stable identity separate from content revision:

```text
DefinitionRevision
  kind and stable identity
  source authority
  language and schema version
  source path and document location
  content digest
  imports and dependency references
  declared target extensions
  authored metadata
```

Default application and normalization produce derived compiler state. They do
not rewrite the source unless an explicit authoring or migration operation is
requested.

### Component contract

A component contract describes portable semantics rather than one framework's
private implementation:

```text
ComponentContract
  namespaced kind and immutable version
  publisher and package provenance
  semantic category and accessibility role
  props, slots, events, actions, and state contract
  constraints and default semantics
  portability class
  target capability requirements
  deprecation and compatibility policy
```

Built-in contracts are versioned parts of the Sigil language distribution.
External contracts remain owned by their publishers. A project registration
can extend or intentionally replace a contract only through an explicit,
namespaced override binding. A same-name file must not silently replace a
built-in semantic contract.

### Design-system binding

A semantic component is realized through a versioned target binding:

```text
ComponentBinding
  component contract reference
  target and renderer profile
  design-system package and version
  runtime import or source reference
  prop, slot, event, and variant mappings
  required styles and peer dependencies
  accessibility and behavior assertions
  supported and unsupported capabilities
  digest and publisher provenance
```

For example, `@hollis-labs/sysop-ui` owns its React components and public APIs.
Sigil owns a binding that says how a portable component maps to that release.
Sigil does not fork the component implementation or silently infer compatibility
from a matching component name.

### Data contract binding

Datasource definitions are UI-side bindings to data owned elsewhere:

```text
DataContractBinding
  stable alias in the UI source
  source authority and schema reference
  pinned schema revision or digest
  query and mutation capabilities
  request and response envelope shapes
  field and relation projections
  pagination, sorting, filtering, and error semantics
  authentication integration reference
  mock and preview policy
```

When an OpenAPI document, JSON Schema, protobuf schema, Go type catalog, or
official client already defines the backend contract, Sigil should consume it
through an adapter rather than establish a competing hand-maintained schema.

A datasource manifest can remain a first-class authored projection where no
external schema is suitable. It is authoritative for what the UI expects, not
for what the backend actually serves.

### UI intermediate representation

The canonical UI IR is the compiler boundary between frontend semantics and
renderers:

```text
UIProgram
  language and IR version
  application and module graph
  resolved pages, routes, overlays, and shells
  resolved component tree and stable node identities
  normalized props, slots, actions, conditions, and shortcuts
  bound data contracts and theme tokens
  required capabilities and effects
  source map and diagnostic locations
  complete dependency and provenance graph
```

Renderers consume the IR, not partially resolved YAML structs, ambient
directories, process globals, or arbitrary repository files. Renderer-specific
extensions are namespaced nodes with an explicit portability consequence.

### Renderer registration

```text
RendererRegistration
  stable target identity and immutable version
  publisher and provenance
  supported language and IR versions
  supported component and capability matrix
  target profiles and options schema
  required toolchain and package dependencies
  output path and ownership contract
  formatter and verifier bindings
  trust class and requested execution capabilities
  digest and compatibility policy
```

Selecting `react-shadcn`, `go-templ`, `catalog`, or another target binds the
exact registration revision. A target name alone is insufficient to reproduce
a historical artifact.

### Artifact set

Renderers return a pure artifact proposal:

```text
ArtifactSet
  compilation identity
  target and profile binding
  ordered files with normalized relative paths
  content, mode, media kind, and digest
  managed, scaffold, or external ownership class
  dependency and toolchain declarations
  source maps and generated-by metadata
  warnings, degradations, and unsupported capabilities
```

Renderers never write directly to the consumer filesystem. The compiler
validates the complete set for collisions, unsafe paths, deterministic order,
and ownership before staging or applying it.

### Compilation receipt

```text
CompilationReceipt
  compiler and IR versions
  source bundle digest
  definition and contract revisions
  renderer and target-profile revision
  effective non-secret configuration
  artifact manifest and digests
  formatter and verifier results
  warnings and accepted degradations
  start, finish, and outcome
  causation, correlation, and trace context
```

The receipt proves what Sigil compiled and what verification it observed. It
does not prove the application was released, deployed, authorized correctly,
or behaved correctly in production.

## Generation authority modes

The phrase “generated code is editable” is incomplete without an ownership
contract. Sigil supports distinct modes rather than blending them.

### Managed projection

In managed mode, the source bundle is authoritative and generated files are
rebuildable projections:

- files carry a generated marker and compilation provenance
- an ownership manifest records every managed path and digest
- manual edits are divergence, not a supported customization mechanism
- regeneration detects modifications before replacing them
- stale managed files may be removed through manifest reconciliation
- handwritten files outside the owned set are untouched
- customization enters through source definitions, component bindings,
  templates, partials, imports, or other declared extension points

Managed output can be version-controlled and reviewed. Version control does not
make hand-editing the projection safe.

### Scaffold transfer

In scaffold mode, Sigil deliberately transfers ownership of the generated
files to the consuming project:

- the receipt records the scaffold generation event
- the consumer may edit the files freely
- Sigil does not later claim those paths as managed output
- re-running a scaffold into existing adopted files requires a distinct merge
  or comparison operation, not silent overwrite
- the source definition may remain useful documentation, but it is no longer a
  complete source of truth for the adopted implementation

Handler stubs are a natural scaffold use case unless their structure is fully
managed around explicit handwritten extension seams.

### Mixed project

A project can contain managed and adopted artifacts when the boundaries are
path- and file-explicit. A generated route table can be managed while custom
page components are adopted; generated interfaces can be managed while
business implementations are handwritten.

One file is not partly managed through conventions such as “edit below this
comment” unless a structured merge format makes that boundary mechanically
enforceable. Prefer separate files, partial methods, interfaces, imports, and
package extension points.

### External package reference

Some output is not copied or generated. Sigil emits imports and configuration
against a versioned external package. The package manager owns materialization;
the design-system publisher owns implementation; Sigil records the binding.

## Compilation pipeline

The desired pipeline is deterministic and side-effect-free until artifact
application:

```text
discover
  -> parse with source locations
  -> resolve versions/imports/registrations
  -> normalize without rewriting source
  -> validate portable semantics
  -> compile canonical UI IR
  -> negotiate target capabilities
  -> render ArtifactSet
  -> validate paths/collisions/ownership
  -> format in a staging root
  -> verify in the declared toolchain
  -> produce diff and receipt
  -> atomically reconcile approved artifacts
```

Preview, dry-run, diff, CI, MCP, and normal generation use this same pipeline.
They differ only in which terminal operation is authorized.

### Discovery

Discovery begins from a resolved project configuration and explicit source
roots. It records every input and rejects ambiguous duplicate identities.
Configured page, theme, datasource, component, provider, and module roots are
actually honored; conventions are defaults, not a parallel hidden authority.

### Parsing and source preservation

Parsing retains source locations, comments, ordering, scalar style, anchors,
and unknown extension nodes where the format promises round-trip preservation.
Typed semantic views can be derived from the YAML node tree.

Read-only compilation never mutates source. Authoring and migration operations
produce an explicit patch or new document revision.

### Resolution

Resolution binds every reference to a stable identity and revision:

- component kinds and design-system mappings
- module pages and shell references
- routes and navigation targets
- datasource aliases and schema sources
- theme inheritance
- custom component and provider sources
- renderer and target profile
- formatters and verifiers

Cycles, ambiguity, missing references, incompatible versions, and namespace
collisions are first-class diagnostics.

### Validation

Validation has layers:

- syntax and language-version validation
- schema validation
- reference and graph validation
- portable semantic validation
- data and action contract validation
- target capability validation
- output and ownership validation
- security and policy validation

Warnings have stable codes and severity. A policy can promote selected warnings
to errors. “Valid YAML” and “portable UI program” are different claims.

### Rendering

Rendering is a pure function of the pinned UI IR, renderer registration,
target profile, and declared build inputs. A renderer cannot inspect arbitrary
repository state, read secrets, access the network, or write files unless its
trust class and capability grant explicitly permit an external tool adapter.

### Formatting and verification

Formatting and verification are declared pipeline adapters. Examples include
`gofmt`, `templ generate`, the Go compiler, TypeScript, framework build tools,
Prettier, and package-specific validators.

The tool, version, configuration, and input scope are recorded. Ambient `PATH`
is not a reproducibility contract. Official CLIs and SDKs are preferred where
they already implement the target format.

`generated`, `formatted`, `type-checked`, `built`, and `verified` are separate
outcomes.

### Application

Artifacts are staged under a validated temporary root. The reconciler checks:

- the output root is explicit, normalized, and safe
- no artifact path is absolute or escapes the root
- symlinks cannot redirect writes outside the root
- artifact paths are unique under platform path semantics
- managed paths do not overlap adopted or unknown files
- current managed files still match their last recorded digests
- the staged set is complete before any replacement

Managed reconciliation is atomic at the strongest boundary supported by the
platform. A failed compilation or formatter never leaves a partially updated
output tree.

Broad recursive deletion is not the generation model. Cleanup removes only
previously managed paths recorded in the ownership manifest unless the operator
explicitly targets a verified, Sigil-exclusive disposable root.

## Lifecycle separation

### Definition lifecycle

```text
authored -> parsed -> resolved -> validated -> compiled
    |         |          |           |
    +---------+----------+-----------+-> invalid
                         +------------> incompatible
                         +------------> deprecated
```

These are observations about a definition revision. Editing creates a new
revision rather than mutating a historical compilation fact.

### Compilation lifecycle

```text
requested -> planned -> rendered -> staged -> verified -> applied
    |           |          |          |          |
    +-----------+----------+----------+----------+-> failed
                                    \-> verification-failed
```

A rendered artifact set is not applied output. Applied output is not a verified
application. A verified application is not a released or deployed application.

### Managed artifact lifecycle

```text
proposed -> applied -> current
              |          |
              |          +-> stale
              |          +-> diverged
              |          +-> orphaned
              +------------> superseded
```

`Stale` means a newer source or compiler binding would produce a different
artifact. `Diverged` means the managed file was modified outside its authority
path. `Orphaned` means its prior owner is missing or its source binding no
longer resolves.

### Scaffold lifecycle

```text
proposed -> transferred -> adopted by consumer
```

After transfer, Sigil does not assign managed lifecycle states to the file.

### Preview lifecycle

```text
source observed -> compiled projection -> served -> invalidated
```

Preview state is derived and disposable. It never becomes the source of UI
truth.

## Portable language and target capabilities

### Honest portability

“One config language, multiple render targets” means that the portable subset
has defined semantics across targets. It does not mean every target silently
approximates every feature.

Each UI program produces a capability requirement set. Each renderer publishes
a support matrix. Compilation selects one of three explicit outcomes per
capability:

- supported with equivalent semantics
- supported through a declared target-specific realization
- unsupported

Lossy degradation requires an explicit policy and appears in the compilation
receipt. Strict mode rejects it.

### Core and extensions

The language contains:

- a portable core for layout, content, forms, data display, navigation,
  overlays, actions, conditions, shortcuts, and design tokens
- versioned semantic component packs
- namespaced target extensions for concepts that are genuinely framework- or
  design-system-specific

Framework-specific concepts such as Next.js route groups, Vite environment
access, React context providers, HTMX attributes, Tailwind class names, and
shadcn imports belong in target profiles or extensions, not the portable core.

An escape hatch must disclose that the source is no longer fully portable.

### Renderer contract shape

Renderers declare granular capabilities rather than implementing one widening
interface filled with no-op methods. Relevant compiler backends may include:

- page and component renderer
- application shell and route renderer
- datasource client generator
- server scaffold generator
- theme and token generator
- provider or dependency wiring generator
- preview renderer
- static catalog renderer

A target profile composes the capabilities it supports. A Go/Templ renderer
does not pretend to support React provider generation by returning no files.

## Component, design-system, and catalog model

### Semantic components

Component schemas define behavior and constraints, not merely prop shapes.
Accessibility roles, keyboard behavior, focus semantics, empty/loading/error
states, controlled versus uncontrolled state, and event meaning are part of the
contract where relevant.

Defaults belong to an identified authority:

- semantic defaults in the component contract
- visual defaults in the design-system binding or theme
- project defaults in project configuration
- target defaults in the renderer profile

Effective values are traceable rather than flattened without provenance.

### Custom components

A custom component registration combines a semantic contract and one or more
target bindings. Source files remain project- or package-owned inputs.

Sigil validates:

- namespaced identity and version
- declared exports and imports
- file containment and digests
- prop, slot, event, and action mapping
- target compatibility
- required dependencies and styles
- ownership mode for copied or referenced source

Copying arbitrary files listed in a schema is not a general filesystem
capability.

### Component packs

Reusable catalogs ship as versioned component packs through the Hollis Labs
plugin SDK or a language-specific package mechanism. A pack can include
contracts, target bindings, examples, theme fragments, verifier rules, and
catalog metadata.

The pack publisher owns those definitions. Sigil resolves, validates,
materializes, caches, and records them. Project-local definitions remain useful
without requiring a central registry.

### Catalog

The component catalog is a static documentation and inspection artifact
generated from registered contracts, examples, variants, target bindings, and
themes. It is not an authoring database or runtime component registry.

Catalog examples are publisher-owned inputs with validation and provenance.
The catalog renderer reuses the same component semantics and target binding as
normal generation so visual documentation cannot silently describe a different
component API.

## Datasources, actions, and runtime security

### Datasource semantics

Datasource capabilities describe what generated UI code may request. They do
not guarantee that the backend supports, authorizes, or correctly implements
the operation.

Request envelopes, response envelopes, errors, pagination, nullability,
filtering, sorting, relations, and mutations are explicit. A single hardcoded
`{data, meta}` REST convention cannot represent every backend contract.

Target renderers generate against the bound contract revision. They do not
invent response unwrapping or infer semantics from endpoint names.

### Handler generation

Generated handlers are either:

- managed adapters around stable user-implemented interfaces, or
- scaffold output transferred to the backend owner

Sigil does not generate placeholder business logic into a file it later
overwrites. Authentication, authorization, validation against current domain
state, transactions, and data access remain backend responsibilities.

### Actions and effects

The UI language can declare intent such as navigate, submit, refresh, emit,
open, close, or request an HTTP operation. The renderer maps that intent into
runtime code.

Sigil does not execute those effects during compilation or preview. Preview
uses inert, mocked, or explicitly sandboxed behavior. URLs, templates, payload
bindings, and methods are validated as code-generation inputs, not treated as
authorization.

### Guards

Page guard metadata can generate route integration points or documentation.
It is not a security boundary. The consuming runtime and backend enforce
identity and authorization. Compilation diagnostics should make this explicit
whenever a guard is only a client-side projection.

### Mock data

Preview and catalog mock data is generated, synthetic, or explicitly supplied
under a non-production policy. Sigil does not load production databases or API
credentials to make previews realistic.

## Theme and design-token boundary

Themes are publisher-owned design definitions. Sigil resolves inheritance,
validates token types and references, and materializes target artifacts.

The target model supports:

- stable token identity and semantic aliases
- typed token values
- modes such as light, dark, high contrast, or product variants
- component token contracts
- deterministic inheritance and override precedence
- cycle detection
- target-specific serialization without target-specific source meaning
- provenance from final token value to defining source

CSS custom properties are one target realization, not the theme authority.
Tailwind configuration is a generated projection. A design-system package may
own the canonical tokens and expose them to Sigil through a versioned adapter.

## Renderer and plugin model

### Extension classes

Sigil can be extended through narrow classes:

- source importers and schema adapters
- component and design-system packs
- render targets and target profiles
- formatter adapters
- verifier and build adapters
- preview fixtures and catalog examples

The Hollis Labs plugin SDK supplies common identity, version, registration,
configuration, capability, lifecycle, health, and diagnostics mechanics.
Sigil defines the UI-compiler-specific contracts.

### Trust and isolation

A renderer or formatter can produce executable source and may invoke external
toolchains. Trust classes are explicit:

- core trusted code shipped with Sigil
- signed portfolio plugin
- declarative pack with no executable code
- isolated external renderer process
- isolated formatter or verifier tool

In-process Go renderers are reserved for trusted code. Third-party executable
extensions run behind an isolation boundary and receive only the UI IR and
declared inputs. They return an artifact set and diagnostics; they do not get
ambient filesystem, network, environment, or credential access.

### Registration and configuration

Plugin configuration contains non-secret values and references. A renderer
registration declares required toolchains and package dependencies; it does not
install arbitrary software during compilation.

Registration does not mean available. Resolution, compatibility, toolchain
materialization, health, and verification are distinct facts.

### Official tools first

Sigil prefers official schemas, SDKs, code generators, formatters, and CLIs
before implementing bespoke equivalents. Adapters preserve Sigil's compilation
receipt and ownership model around those tools.

Examples include official OpenAPI tooling, framework compilers, `gofmt`, Templ,
TypeScript, Prettier, and package managers. An external tool's success is
verification evidence, not proof of semantic UI equivalence.

## Determinism and reproducibility

A reproducible compilation binds:

- compiler version and build digest
- UI language and IR versions
- complete source bundle digest
- component, design-system, and schema package revisions
- renderer and target profile revision
- formatter and verifier versions
- effective non-secret configuration
- declared environment and platform dimensions

The compiler guarantees:

- sorted and stable traversal where source order is not semantic
- stable diagnostic ordering
- no timestamps, random IDs, host paths, or map iteration in generated content
  unless explicitly declared
- normalized line endings and modes
- stable output paths
- no network lookups during the hermetic compilation phase
- content-addressable cache keys derived from all relevant inputs

Two equivalent builds under the same declared platform contract produce the
same artifact digests. Target toolchains that cannot guarantee byte identity
must declare the weaker reproducibility guarantee.

## Output safety and reconciliation

### Ownership manifest

Every managed output root contains or is associated with an ownership manifest:

```text
OutputOwnershipManifest
  project and output-root identity
  generation authority mode
  compiler and renderer binding
  last source bundle and compilation identity
  managed paths and digests
  adopted and protected path declarations
  target platform path semantics
```

The manifest is operational metadata. UI definitions remain in the source
bundle.

### Path safety

All source and output paths are resolved against declared roots. Sigil rejects:

- absolute artifact paths
- `..` escapes
- symlink or junction escapes
- platform-reserved names
- case-folding collisions
- duplicate normalized paths
- output roots equal to a repository root, home directory, filesystem root, or
  another broad target unless a specific safe policy permits it
- custom source references outside approved roots

### Atomicity

The complete artifact set is rendered, formatted, validated, and staged before
application. A failure leaves the prior applied set intact.

Source authoring writes use expected digests and atomic replacement. Concurrent
CLI, MCP, editor, and Studio changes either compose through structured patches
or fail with a conflict; last-writer-wins file replacement is not the model.

### Drift

Sigil can explain:

- source drift: definitions changed since the last compilation
- dependency drift: a component, schema, renderer, or toolchain binding changed
- generated drift: a managed artifact was edited
- stale output: a previously managed file is no longer proposed
- verification drift: the same artifact no longer passes its verifier contract

Diffs compare semantic definitions, compiled IR, and artifacts as distinct
layers.

## Authoring and migration

### Authoring service

CLI, MCP, and a possible Studio use one source-authoring service:

```text
GetDefinition
ListDefinitions
ValidateDefinition
CreateDefinition
ApplyDefinitionPatch
MoveDefinition
DeleteDefinition
MigrateDefinition
Compile
DiffCompilation
InspectCapability
```

Mutations are rooted to one declared project and require an expected source
revision. Operations return the new digest, diagnostics, and semantic diff.

### Round-trip preservation

Authoring and migration preserve comments, ordering, formatting style, unknown
extensions, and source locations through the YAML node representation where
the source contract promises it. A whole-document marshal is suitable for a
new file but not for editing an established hand-authored document.

### Migrations

Language migrations are versioned, deterministic transformations:

- the original revision remains recoverable through version control or an
  explicit backup artifact
- dry-run and semantic diff are first-class
- each transformation reports its reason and affected source locations
- migrations are idempotent
- no default is written merely because the in-memory compiler applied it
- stable node identity is preserved where possible
- unsupported lossy transformations stop for user direction

Migration changes authored source; compilation does not.

### Import and export

JSON export is a projection of the semantic model unless it explicitly claims
lossless source round-trip. Import validates authority and collision policy.
Converting YAML to JSON and back must not silently claim to preserve comments,
anchors, scalar style, or ordering.

## Interfaces and parity

### Public Go API

Sigil exposes a stable public compiler and authoring API rather than requiring
consumers such as a sibling Studio to import `internal/` packages. The API uses
versioned DTOs and narrow interfaces for source access, registry resolution,
rendering, artifact application, and verification.

The CLI and MCP server are adapters over this API. There is one validation and
compilation implementation.

### CLI

The CLI is the primary human, scripting, and build-pipeline surface. Commands
separate read-only inspection, source mutation, compilation, artifact
application, and verification.

Destructive options show the resolved root and affected managed paths. A
generic `--clean` never translates directly into recursively deleting an
unverified output directory.

Exit status, structured output, diagnostic codes, and machine-readable receipts
make the CLI suitable for CI without parsing presentation text.

### MCP

MCP is the primary portable agent-facing authoring surface. Sigil uses the
official Go SDK and current protocol semantics rather than maintaining its own
general JSON-RPC implementation.

MCP tools and resources project the same application service as the CLI:

- discovery exposes the effective project, component, datasource, theme,
  renderer, and capability catalog
- mutations use expected revisions and structured patches
- validation and compilation return typed diagnostics and receipts
- generation does not gain broader filesystem authority merely because an
  agent invoked it
- tool availability and docs are generated from or verified against one
  registry

A stdio process may be created per agent client because project files are the
authority and Sigil has no required daemon state. Concurrent processes still
coordinate through revision checks and filesystem-safe writes.

Tether may proxy or catalog Sigil's MCP server when composed. Direct local MCP
remains fully supported.

### Preview server

`sigil serve` is a disposable local development process that:

- watches a declared source bundle
- invokes the same compiler pipeline into an in-memory or temporary preview
  projection
- serves static preview and catalog artifacts
- publishes reload and diagnostic events
- exposes liveness and compilation readiness

It binds loopback by default. Network exposure requires explicit bind,
authentication, origin, and content-security policy. Preview actions are inert
or sandboxed, and source content is treated as untrusted HTML input.

The preview server owns no durable business state and does not need to run as a
portfolio daemon. If an operator wants an always-on preview, Cerberus can
materialize an OS service definition and the OS supervises it. That deployment
choice does not change Sigil's compiler boundary.

### Studio

If Sigil Studio or another visual authoring application exists, it is a sibling
product. It owns canvases, windows, authoring sessions, selection state, undo,
collaboration history, and visual editing UX.

Studio consumes Sigil's public compiler and authoring contracts. Sigil does not
depend on Studio, absorb Wails or tldraw, or create a second document store.
Compose-mode definitions remain project-owned files. Scene-mode artifacts
remain Studio-owned until explicitly translated or referenced.

The lightweight combination of catalog, preview, MCP, Tangent, and chat may
cover many review and iteration flows. Whether Studio is warranted is a product
decision outside Sigil core; the compiler boundary remains the same either way.

## Process and persistence model

Sigil is primarily a stateless command process:

- project files are the authored source of truth
- external packages and registries are attached definition sources
- generated ownership manifests and compilation receipts are operational
  artifacts
- build caches and preview caches are derived and disposable
- no database is required for core compilation

CLI and stdio MCP invocations load one source snapshot, perform an operation,
emit results, and exit. The preview server retains only live connections,
watch state, diagnostics, and derived compilation caches.

If remote registries or expensive compilation justify a shared cache, the cache
is content-addressed and rebuildable. It does not become a mutable source of UI
truth.

Sigil does not self-daemonize, maintain PID files, or manage OS services.
Long-running preview processes run in the foreground and support graceful
shutdown. Cerberus may install the binary and materialize an OS service; the OS
or runtime owns supervision.

## Configuration ownership

Configuration is external, deterministic, and non-secret. It describes:

- project and source roots
- language and schema versions
- registered component, design-system, renderer, and verifier packages
- target profiles and output authority modes
- theme and datasource bindings
- output roots and ownership policy
- formatter and verifier toolchain references
- preview listener policy
- cache and receipt locations
- telemetry and redaction policy

Source definitions and runtime deployment configuration remain separate.
Application API base URLs may be generated as runtime environment-variable
references; environment-specific values do not belong in the portable UI
definition.

Precedence among project config, target profile, CLI flags, and environment
bindings is explicit and inspectable. Every declared setting is either honored
or rejected; unused configuration is a diagnostic rather than silent drift.

Effective configuration is recorded in the compilation receipt without secret
values or machine-specific noise.

## Secrets and supply-chain security

Sigil does not own API keys, OAuth tokens, registry passwords, package-manager
credentials, signing keys, or deployment secrets.

Definitions may name environment variables, credential references, or runtime
injection points. They never carry actual secret values. Generated examples and
preview mock data contain no production credentials.

Remote package or registry access is performed by an external credential-aware
dependency mechanism or a narrowly scoped resolver grant. Secret material is
not passed into renderers, templates, generated source, diagnostics, receipts,
logs, or MCP resources.

Supply-chain policy verifies:

- compiler and plugin provenance
- package versions, digests, and signatures where available
- renderer and component-pack compatibility
- external toolchain identity
- generated dependency declarations
- executable extension trust class
- artifact provenance suitable for the consuming build system

Templates, component source, generated code, and preview HTML are untrusted
inputs until validated under their relevant execution boundary.

## Events, diagnostics, and observability

Sigil emits structured events and diagnostics such as:

- source bundle discovered
- definition parsed, defaulted, migrated, resolved, or rejected
- component or renderer registration resolved or incompatible
- target capability supported, degraded, or missing
- compilation planned, rendered, staged, verified, applied, or failed
- managed artifact current, stale, diverged, superseded, or orphaned
- preview compilation invalidated or recovered
- authoring conflict detected
- external formatter or verifier invoked and completed

Compiler diagnostics include stable code, severity, message, source range,
definition identity, target, related locations, and optional corrective hint.
Human and machine interfaces use the same diagnostic model.

OpenTelemetry spans can correlate source resolution, validation, IR
compilation, rendering, formatting, verification, and application. Useful
metrics include:

- compilation latency by stage and target
- source, IR, and artifact cache hit rates
- diagnostic counts by code and severity
- unsupported and degraded capability counts
- renderer and verifier failures
- generated file count and byte size
- managed drift and collision counts
- preview rebuild and connected-client counts

Telemetry excludes source bodies, generated code, file contents, user data,
secrets, and sensitive paths by default. Logs are structured event streams to
stdout and stderr. Compilation receipts are artifacts, not log files.

## Portfolio composition

### Tangent

Tangent renders typed interactions at runtime. Sigil generates source at build
time. Tangent has no Sigil runtime dependency.

Sigil may scaffold a Tangent renderer component, generate schema bindings, or
produce a static design catalog. Tangent or another caller may present a Sigil
semantic diff, preview, or compilation receipt for human review. Tangent owns
that interaction; Sigil owns the compilation fact; the UI publisher owns the
definition.

### Nanite

Nanite owns agents, sessions, conversations, context, and runtime interaction
modes. An agent in Nanite can discover components, edit project definitions,
validate, compile, and inspect artifacts through Sigil's MCP or public API.

Nanite does not become the UI definition authority merely because its agent
authored a change. The project repository and its review policy remain
authoritative.

### Tether

Tether is optional. It may route Sigil's MCP surface, provide caller identity,
or expose Sigil through a portfolio catalog. Sigil remains directly usable as
a CLI and stdio MCP server.

Tether owns gateway, messaging, and session policy. Sigil owns compiler and
authoring semantics. Neither silently owns the other's settings.

### Cerberus and Coder

Cerberus may install a pinned Sigil binary, materialize a foreground or
OS-managed preview-service definition, provide attached toolchains, and observe
health. The OS or runtime supervises any long-running preview process.

Cerberus does not own UI definitions or generated source. A Cerberus-managed
Coder Workspace may provide the filesystem and build environment in which
Sigil runs. Sigil consumes the checked-out project and declared toolchain; it
does not provision, lease, or manage the Workspace.

### Hadron

Hadron may orchestrate a reusable UI build, verification, review, and publish
workflow. It owns the workflow run and gates. Sigil owns each compilation and
receipt. The consuming build system owns release acceptance.

Sigil itself does not acquire workflow retries, scheduling, or human-gate
semantics.

### Torque

Torque may track UI work, request a compilation, attach generated diffs or
receipts as evidence, and apply acceptance policy. Sigil does not create,
schedule, lease, or complete tasks.

### Tesseract

Tesseract may index explicit pointers to UI definitions, component contracts,
architecture decisions, or compilation receipts for recall. Sigil does not
automatically promote project source or generated code into portfolio memory.

### Fragments Engine and Loom

Fragments Engine may deliver source material or design references to an
authoring workflow. Loom owns generated content, templates, publications, and
output quality. Sigil owns UI compilation only.

A Loom page or Fragments Engine record can be displayed by a Sigil-generated
application through an external data contract. Sigil does not become the
content authority.

### Design-system and shared UI libraries

`sysop-ui`, shadcn, Templ components, and other UI packages own their runtime
implementations and public APIs. Sigil consumes versioned adapters and emits
imports or managed bindings. Renderer dogfooding against real applications is
valuable evidence, but generated parity is measured rather than assumed.

## Twelve-factor and Go operating model

Sigil follows the portfolio's adapted twelve-factor direction:

- One version-controlled codebase produces versioned Sigil binaries and plugin
  packages for many environments.
- Go dependencies, frontend target dependencies, formatters, and verifiers are
  explicitly declared and isolated. No build relies on unexplained global
  tools.
- Configuration is external and contains no secret values. Environment
  variables bind deployment-specific values and credential references.
- Source repositories, schema registries, package registries, caches, preview
  listeners, and verification toolchains are attached resources.
- Build creates the immutable Sigil compiler; release binds it to plugin and
  schema packages; run performs a compilation or serves a disposable preview.
  The generated application has its own independent build, release, and run
  lifecycle.
- CLI and MCP compiler processes are stateless. Durable definitions and
  artifacts live in version control or declared attached stores.
- The optional preview service embeds its HTTP server and binds an explicitly
  configured local port.
- Concurrency scales through independent hermetic compilations and safe output
  reconciliation, not shared process-global registries or competing writes.
- Startup is fast, shutdown is graceful, and interruption before artifact
  application leaves the prior output intact.
- Development and CI use the same compiler, renderer, definition, and verifier
  contracts as release builds.
- Logs go to stdout and stderr; source history remains in revision control and
  compilation receipts remain structured artifacts.
- Schema migrations, cache repair, validation, compilation, verification,
  catalog generation, import, and export run as one-off processes from the same
  release.

Go-specific conventions follow the broader Hollis Labs engineering direction:

- domain and compiler packages do not depend on Cobra, MCP, HTTP, filesystem
  layout, or a particular renderer
- transports and storage adapters depend inward on public application contracts
- constructors validate required dependencies and fail fast
- contexts carry cancellation and deadlines, not optional dependencies
- interfaces are consumer-owned and intentionally narrow
- errors and diagnostics remain typed across adapters
- process-global mutable registries are avoided or frozen during composition
- maps are sorted before they influence output or diagnostics
- official Go libraries and protocol SDKs are preferred
- generated code records its source and is reproducibly checked

## Current strengths to preserve

The implementation already contains strong architectural choices:

- an explicit build-time code-generation boundary
- no generated-application runtime dependency on Sigil
- a declarative, human- and agent-readable YAML source format
- a system-agnostic component and theme ambition
- multiple framework-native render targets
- a renderer contract returning `[]OutputFile` instead of writing files directly
- validation before page generation and before MCP page writes
- embedded built-in component schemas
- project-local custom component contracts and source files
- target-native React, Go/Templ, HTMX, TypeScript, CSS, and Tailwind output
- dry-run support
- semantic config diff, migration, import, export, schema export, and doctor
  surfaces
- a small Go dependency graph
- a disposable live preview server with reload support
- a portable agent-facing MCP surface
- a static-preview foundation suitable for a component catalog
- generated-source readability and reviewability
- per-module provider scoping rather than app-global leakage
- preservation of handwritten demo hooks through explicit generation scripts
- direct dogfooding against Stack Explorer, Clockwork, and Sigil's own UI
- honest prior decisions not to replace hand-composed screens when generated
  output lacked runtime parity
- the existing decision that a visual Studio is a sibling rather than part of
  compiler core

The target should clarify these strengths, not turn Sigil into a dynamic UI
database, hosted editor, runtime framework, or general application generator.

## Architectural tensions to resolve

These are target-state design questions exposed by the current implementation,
not a delivery backlog.

### Config authority conflicts with editable generated code

Current architecture says the YAML fully describes the UI while also describing
generated code as editable. Regeneration cannot preserve both claims without an
explicit managed, scaffold, or mixed ownership mode.

### One config, multiple targets overstates semantic parity

Go/Templ and React renderers do not implement the same application shell, API,
provider, component, and action capabilities. Returning no files from an
unsupported renderer method hides incompatibility instead of reporting it.

### Framework concepts leak into the portable model

`target_mode`, Next.js route groups, Vite and Next environment access, React
context providers, source `.tsx` files, Tailwind `className`, and shadcn-oriented
variants appear in general app and component configuration. Target profiles and
namespaced extensions need to contain these concepts.

### The renderer interface is one widening capability bundle

Every renderer must implement page, theme, datasource, shared component,
layout, API client, and provider methods even when the target does not support
them. This encourages silent no-op behavior and makes preview or catalog targets
pretend to be full application generators.

### There is no canonical compiled IR

Renderers receive mutable config structs, a live component registry, maps,
project configuration, and the `.sigil` directory. Resolution and target logic
are spread between the engine and renderers, and renderers can read additional
ambient files.

### Compiler composition uses process-global registration

Renderers register through package `init` into a mutable global map. Available
renderer order is nondeterministic, duplicate names silently replace prior
registrations, and different application embeddings cannot compose isolated
registries.

### Custom schemas silently replace built-ins

Project custom schemas can overwrite a built-in type by name. There is no
namespace, version, publisher, compatibility, or explicit override decision.

### Custom source references have broad filesystem meaning

Custom component and provider source paths are joined to the Sigil directory
and copied into renderer-selected paths. Containment, symlink escape, normalized
output collision, and target-specific ownership need one host-enforced policy.

### Clean generation recursively removes the output root

`--clean` currently maps to `os.RemoveAll` on the supplied output path. The
compiler needs resolved-root safety and manifest-based reconciliation rather
than broad deletion.

### Artifact writes are incremental and non-atomic

The engine writes files one at a time after optional deletion. A late failure
can leave a partially generated tree, and duplicate artifact paths can overwrite
one another without a complete preflight.

### Generated ownership is maintained by shell conventions

Demo targets use temporary directories, `rsync`, `--ignore-existing`, selective
copying, `sed`, and manually protected hooks to express which files Sigil owns.
These scripts reveal the correct need but make ownership external and
target-specific rather than a compiler contract.

### Stale and modified generated files are not modeled

There is no output ownership manifest, source-to-artifact digest binding, or
first-class current, stale, diverged, adopted, and orphaned vocabulary.

### Generation is not fully deterministic

Several emitted structures and CLI displays originate from Go maps. Known map
iteration nondeterminism has already caused regeneration diff churn and blocks
reliable byte-for-byte parity gates.

### Project configuration is partly ignored

Project config declares default output, clean policy, component, datasource,
and theme directories, while core generation paths and CLI defaults hardcode
several `.sigil` and `internal/ui` conventions. Declared but unused settings
create two apparent authorities.

### Theme inheritance lacks a complete graph contract

Theme inheritance recursively loads and mutates a base object without an
explicit cycle, identity, provenance, or multi-mode model. The resulting theme
can obscure whether values came from the base or child.

### Datasource assumptions are renderer policy disguised as schema

The documented `{data, meta}` REST envelope and generated SWR hooks do not match
all real application APIs. Sigil's own sysop dogfood exposed that single-key
response envelopes break generated pages. Response shape must be a bound data
contract, not a hardcoded renderer assumption.

### Handler stubs blur managed and handwritten ownership

Generated Go handlers instruct developers to fill in query logic while also
warning them not to edit structure. Regeneration has no mechanical way to
preserve that mixed ownership safely.

### Client-side guards can be mistaken for authorization

Page metadata can declare `auth` or `admin` guards, but Sigil does not own the
runtime identity or backend enforcement. Generated UI affordances must not imply
a security guarantee.

### MCP documentation and implementation have drifted

Documentation advertises generation, preview, theme, and datasource operations
beyond the nine tools registered by the current server. Tool docs, discovery,
and implementation do not share one authoritative registry.

### MCP is a bespoke protocol implementation

The current stdio server implements an older MCP JSON-RPC surface directly.
Using the official SDK would reduce protocol ownership and align capabilities,
cancellation, structured content, resources, and version negotiation with the
standard.

### MCP mutations lack revision and path authority

Page updates replace whole files without expected digests or atomic writes.
Datasource creation writes unvalidated YAML. Read and validate paths can reach
caller-supplied filesystem locations beyond the project root. Concurrent agents
and editors have last-writer-wins behavior.

### Preview binds beyond loopback

The live server listens on `:<port>`, while its output and source inspection are
intended for local development. It has no authentication, origin policy,
readiness model, or generalized untrusted-content boundary.

### Preview and generation do not share one compiled truth

The live server loads pages, themes, and datasources through its own partial
paths and invokes the preview renderer directly. It can accept inputs or show
semantics that differ from normal validation and generation.

### Preview is deployed as a persistent service

Cerberus currently materializes `sigil serve` as a launchd-backed service. That
is a valid operator convenience, but a repository-watching dev preview is not a
required Sigil daemon and should not be mistaken for a production runtime.

### Built-in registry failure can silently produce an empty compiler

The default registry treats embedded-schema load failure as an empty registry.
A compiler distribution missing its core language should fail readiness rather
than defer the error to unrelated validation messages.

### Component contracts validate shape more than behavior

Schemas capture props, actions, slots, and shortcuts, but accessibility, focus,
state, data envelope, error, and target equivalence semantics are incomplete.
Forty-nine registered names do not by themselves prove forty-nine portable
components.

### Version and repository identity are inconsistent

The project is a Hollis Labs application while the Go module, install examples,
linker path, and imports still use `github.com/hollis-labs/sigil`. Release and
package authority should have one canonical identity or an explicit transition
contract.

### Build and release evidence is incomplete

The repository has Make-based local build, test, and vet commands but no current
CI workflow in the project metadata. Renderer output verification is uneven and
compilation receipts do not yet capture the build toolchain and results.

### Architecture and release documentation lag current behavior

Docs describe an earlier component count and renderer shape, while the current
code contains multi-module apps, SPA mode, custom providers, and extensive
dogfood-specific behavior. Stable architecture, current capability, and roadmap
material need distinct authority.

### Dogfood-specific aliases live in the generic engine

The compiler contains a hardcoded `chart -> sigil-chart` custom-component alias
and renderer behavior shaped by particular demos. Project and design-system
bindings should express these mappings without adding product-specific cases to
compiler core.

### Studio-oriented editing could pull runtime UX into core

Round-trip patching, catalogs, previews, and MCP authoring are appropriate
compiler surfaces. Canvas state, windows, collaboration, turn history, and
visual editing are a sibling application's domain. The boundary needs to remain
visible as authoring ambitions grow.

## Boundary guidance

When deciding whether a capability belongs in Sigil, use these tests.

It belongs in Sigil when it primarily:

- parses, migrates, resolves, or validates UI definitions
- defines portable UI compiler semantics
- compiles source into a canonical UI IR
- negotiates renderer and target capabilities
- produces deterministic framework-native artifacts
- manages explicit generated-file ownership and drift
- validates component, theme, datasource, and action bindings
- formats or verifies generated artifacts through declared adapters
- provides project-rooted CLI, MCP, Go API, preview, or catalog projections of
  the same compiler service
- operates Sigil's own registries, caches, diagnostics, and receipts

It probably belongs elsewhere when it primarily:

- renders dynamic business data in production
- stores application or user state
- authenticates users or enforces authorization
- executes application actions or backend mutations
- owns a design-system implementation
- owns backend schemas, databases, or API behavior
- provides a visual canvas, desktop shell, collaborative document, or undo
  history
- manages agent sessions, prompts, conversations, or model calls
- runs general workflows or coordinates tasks
- brokers general MCP, LLM, or messaging traffic
- provisions infrastructure or Coder Workspaces
- stores or rotates secrets
- owns generated content or publication semantics

If the answer is mixed, keep semantic authority in the owning project and
compose through a versioned `SourceBundle`, `ComponentContract`,
`DataContractBinding`, `UIProgram`, `RendererRegistration`, `ArtifactSet`,
`CompilationReceipt`, or narrow adapter.

## Questions for the next architecture session

1. What is the canonical UI language boundary, and which current app, provider,
   route, action, and class-name concepts move into target extensions?
2. What exact managed, scaffold, external-reference, and mixed ownership modes
   does Sigil support, and how are they declared per artifact class?
3. What ownership manifest and reconciliation semantics replace broad clean,
   overwrite, `rsync`, and ignore-existing conventions?
4. What canonical UI IR preserves source locations, stable node identities,
   effective defaults, capabilities, data bindings, and extension nodes?
5. Which component semantics form the portable core, and what evidence is
   required before claiming equivalent support across targets?
6. How are component contracts namespaced, versioned, published, overridden,
   deprecated, and bound to design-system implementations?
7. Which existing built-ins are semantic components versus shadcn- or
   application-shaped conveniences?
8. What granular renderer capability interfaces replace the current single
   widening `Renderer` contract?
9. Which renderer and component-pack extensions may run in process, and what
   isolation contract applies to third-party executable plugins?
10. What datasource binding model supports OpenAPI, JSON Schema, Go types,
    official clients, paginated envelopes, single-key envelopes, bare arrays,
    streams, and custom adapters without hardcoded renderer assumptions?
11. Which generated backend artifacts are managed adapters, and which are
    one-time scaffolds transferred to the consumer?
12. What source-preserving YAML document model supports structured patches,
    migrations, comments, ordering, unknown extensions, and concurrent agent
    edits?
13. What deterministic compilation receipt and artifact provenance format can
    CI, Cerberus, Torque, Hadron, and consuming repositories understand?
14. Which formatter, build, type-check, accessibility, and parity verifiers are
    part of each target profile, and which results are required versus advisory?
15. What public Go API is stable enough for Studio and other embeddings without
    exposing internal compiler packages?
16. Which MCP authoring operations are required, and how are project root,
    caller identity, expected revisions, atomicity, and tool discovery enforced?
17. Should `sigil serve` remain only an on-demand foreground dev process, or is
    an OS-managed always-on preview a supported deployment profile with explicit
    readiness and security constraints?
18. What target and design-system matrix should the catalog visualize so it
    documents actual renderer capability rather than an abstract schema alone?
19. What is the canonical module and package identity for Sigil within Hollis
    Labs, and what compatibility commitment applies to the current import path?
20. What evidence would justify a sibling Studio beyond catalog, preview, MCP,
    Tangent review surfaces, and existing chat-based authoring?

<script>
  let activeTab = 'overview'; // 'overview', 'topology', 'matrix', 'cicd'

  const serviceSpecs = [
    {
      name: 'Nginx Reverse Proxy',
      container: 'songmoodboard_nginx',
      image: 'nginx:1.27-alpine',
      hostPort: '1122',
      internalPort: '80',
      isolation: 'Public Host Edge',
      protocol: 'HTTP/1.1 Reverse Proxy, Gzip',
      role: 'Single public gateway. Forwards /api to Go API and /* to SvelteKit SSR.'
    },
    {
      name: 'Go REST API Server',
      container: 'songmoodboard_backend',
      image: 'golang:1.23 Alpine Multi-Stage (~15MB)',
      hostPort: 'None (Blocked)',
      internalPort: '8080',
      isolation: 'Private (song-moodboard-net)',
      protocol: 'REST / JSON, MinIO S3 SDK, pgx',
      role: 'Static Go binary (-s -w). Runs as unprivileged appuser (UID 1001).'
    },
    {
      name: 'SvelteKit SSR Frontend',
      container: 'songmoodboard_frontend',
      image: 'oven/bun:1-alpine',
      hostPort: 'None (Blocked)',
      internalPort: '3000',
      isolation: 'Private (song-moodboard-net)',
      protocol: 'HTTP Node Adapter, TanStack Query',
      role: 'Server-side rendered universal frontend executed natively with Bun.'
    },
    {
      name: 'PostgreSQL Database',
      container: 'songmoodboard_postgres',
      image: 'postgres:16-alpine',
      hostPort: 'None (Blocked)',
      internalPort: '5432',
      isolation: 'Private (song-moodboard-net)',
      protocol: 'TCP / PostgreSQL Wire Protocol',
      role: 'ACID relational store. Auto-migrates users, tracks, and moodboard items.'
    },
    {
      name: 'Redis In-Memory Cache',
      container: 'songmoodboard_redis',
      image: 'redis:7-alpine',
      hostPort: 'None (Blocked)',
      internalPort: '6379',
      isolation: 'Private (song-moodboard-net)',
      protocol: 'RESP (REdis Serialization Protocol)',
      role: 'High-throughput catalog caching and JWT token blacklist with automatic TTL.'
    },
    {
      name: 'SeaweedFS Object Store',
      container: 'songmoodboard_seaweedfs',
      image: 'chrislusf/seaweedfs',
      hostPort: 'None (Blocked)',
      internalPort: '8333',
      isolation: 'Private (song-moodboard-net)',
      protocol: 'S3-Compatible HTTP REST API',
      role: 'Lightweight distributed blob storage for audio files and high-res photography.'
    }
  ];

  const pipelineStages = [
    {
      id: '01',
      title: 'Backend Verification',
      command: 'go vet ./... && go test -race ./...',
      desc: 'Performs static analysis and concurrency race condition checking across all Go packages.'
    },
    {
      id: '02',
      title: 'Frontend SSR Build Check',
      command: 'bun install --frozen-lockfile && bun run build',
      desc: 'Validates SvelteKit compilation, adapter-node SSR generation, and dependency tree.'
    },
    {
      id: '03',
      title: 'Container Layer Compilation',
      command: 'docker compose build --parallel',
      desc: 'Builds stripped multi-stage Alpine images with cache mounting for reproducible builds.'
    },
    {
      id: '04',
      title: 'Port Isolation Smoke Test',
      command: 'curl -f http://localhost:1122/api/health && ! nc -z localhost 8080',
      desc: 'Asserts HTTP 200 on port 1122 and validates that internal ports (8080, 5432, 6379, 8333) remain closed to the host.'
    }
  ];
</script>

<div class="max-w-[1200px] mx-auto px-6 sm:px-10 py-12 text-neutral-900 dark:text-neutral-100 font-sans animate-fade-in">
  <!-- Top Technical Header (Linear / Vercel style) -->
  <div class="border-b border-neutral-200 dark:border-neutral-800 pb-8 mb-8">
    <div class="text-xs text-neutral-500 dark:text-neutral-400 mb-2 font-medium">
      System specification &bull; Revision 2.4
    </div>

    <h1 class="text-2xl sm:text-4xl font-semibold tracking-tight text-neutral-900 dark:text-neutral-100">
      System Architecture & Topology
    </h1>
    <p class="text-sm text-neutral-600 dark:text-neutral-400 mt-2 max-w-2xl leading-relaxed">
      Complete structural specification of the Song Moodboard cluster. Outlines the single-ingress security boundary on port 1122, container isolation rules, and multi-stage pipelines.
    </p>

    <!-- Navigation Tabs -->
    <div class="flex items-center gap-1 mt-6 border-b border-neutral-200 dark:border-neutral-800 -mb-8">
      {#each [
        { id: 'overview', label: 'Overview & Ingress' },
        { id: 'topology', label: 'Topology Diagram' },
        { id: 'matrix', label: 'Service Matrix' },
        { id: 'cicd', label: 'CI/CD Pipeline' }
      ] as tab}
        <button
          type="button"
          on:click={() => activeTab = tab.id}
          class="px-4 py-2.5 text-xs font-medium border-b-2 transition-colors {activeTab === tab.id ? 'border-neutral-900 dark:border-neutral-100 text-neutral-900 dark:text-neutral-100' : 'border-transparent text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-200'}"
        >
          {tab.label}
        </button>
      {/each}
    </div>
  </div>

  {#if activeTab === 'overview'}
    <!-- Tab 1: Overview & Ingress Boundary -->
    <div class="space-y-10">
      <!-- Ingress Boundary Specification -->
      <div>
        <h2 class="text-base font-semibold text-neutral-900 dark:text-neutral-100 mb-1">Host Ingress Security Boundary</h2>
        <p class="text-xs text-neutral-500 mb-4">Port isolation rules enforced between host machine and internal container subnet.</p>

        <div class="border border-neutral-200 dark:border-neutral-800 rounded-md overflow-hidden bg-white dark:bg-[#0A0A0A]">
          <div class="grid grid-cols-1 md:grid-cols-3 divide-y md:divide-y-0 md:divide-x divide-neutral-200 dark:divide-neutral-800 text-xs">
            <div class="p-5">
              <span class="text-xs text-neutral-500 font-medium block mb-1">Host Public Port</span>
              <div class="text-xl font-mono font-semibold text-neutral-900 dark:text-neutral-100">1122 / TCP</div>
              <p class="text-neutral-600 dark:text-neutral-400 mt-2 leading-relaxed">
                The only exposed entrypoint. Nginx handles SSL termination, compression, and request multiplexing.
              </p>
            </div>

            <div class="p-5">
              <span class="text-xs text-neutral-500 font-medium block mb-1">Internal Docker Network</span>
              <div class="text-xl font-mono font-semibold text-neutral-900 dark:text-neutral-100">song-moodboard-net</div>
              <p class="text-neutral-600 dark:text-neutral-400 mt-2 leading-relaxed">
                Bridge network with automatic DNS resolution. Services communicate via service names (e.g. <code class="font-mono text-[11px] bg-neutral-100 dark:bg-neutral-800 px-1 py-0.5 rounded">backend:8080</code>).
              </p>
            </div>

            <div class="p-5">
              <span class="text-xs text-neutral-500 font-medium block mb-1">Firewall & Zero Exposure</span>
              <div class="text-xl font-mono font-semibold text-emerald-600 dark:text-emerald-400">Strictly Isolated</div>
              <p class="text-neutral-600 dark:text-neutral-400 mt-2 leading-relaxed">
                PostgreSQL (5432), Redis (6379), and SeaweedFS (8333) have zero bindings to host loopback or interfaces.
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Execution Path Specifications -->
      <div>
        <h2 class="text-base font-semibold text-neutral-900 dark:text-neutral-100 mb-1">Request Execution Paths</h2>
        <p class="text-xs text-neutral-500 mb-4">Step-by-step transaction lifecycle for media streaming, persistence, and auth.</p>

        <div class="border border-neutral-200 dark:border-neutral-800 rounded-md divide-y divide-neutral-200 dark:divide-neutral-800 bg-white dark:bg-[#0A0A0A] text-xs">
          <!-- Flow 1 -->
          <div class="p-5 flex flex-col md:flex-row md:items-start justify-between gap-4">
            <div class="md:w-1/3">
              <span class="text-xs text-neutral-500 font-medium">Route 01</span>
              <h3 class="text-sm font-semibold text-neutral-900 dark:text-neutral-100 mt-0.5">Multipart Media Upload</h3>
              <p class="text-neutral-500 mt-1 font-mono text-[11px]">POST /api/tracks</p>
            </div>
            <div class="flex-1 space-y-1 text-neutral-600 dark:text-neutral-400 leading-relaxed">
              <p>1. Client sends multipart payload to port 1122.</p>
              <p>2. Nginx forwards stream directly to <code class="font-mono bg-neutral-100 dark:bg-neutral-800 px-1 py-0.5 rounded">backend:8080</code> without disk buffering.</p>
              <p>3. Go server verifies JWT session and streams audio/photo directly into SeaweedFS S3 bucket.</p>
              <p>4. PostgreSQL transaction saves metadata, while Redis invalidates catalog cache key.</p>
            </div>
          </div>

          <!-- Flow 2 -->
          <div class="p-5 flex flex-col md:flex-row md:items-start justify-between gap-4">
            <div class="md:w-1/3">
              <span class="text-xs text-neutral-500 font-medium">Route 02</span>
              <h3 class="text-sm font-semibold text-neutral-900 dark:text-neutral-100 mt-0.5">Low-Latency Audio Streaming</h3>
              <p class="text-neutral-500 mt-1 font-mono text-[11px]">GET /api/storage/audio/*</p>
            </div>
            <div class="flex-1 space-y-1 text-neutral-600 dark:text-neutral-400 leading-relaxed">
              <p>1. Browser audio player queries stream endpoint.</p>
              <p>2. Go backend checks Redis cache for track header details.</p>
              <p>3. Media byte-range requests are proxied directly to the client with HTTP 206 Partial Content.</p>
            </div>
          </div>

          <!-- Flow 3 -->
          <div class="p-5 flex flex-col md:flex-row md:items-start justify-between gap-4">
            <div class="md:w-1/3">
              <span class="text-xs text-neutral-500 font-medium">Route 03</span>
              <h3 class="text-sm font-semibold text-neutral-900 dark:text-neutral-100 mt-0.5">Stateless Authentication</h3>
              <p class="text-neutral-500 mt-1 font-mono text-[11px]">POST /api/auth/login</p>
            </div>
            <div class="flex-1 space-y-1 text-neutral-600 dark:text-neutral-400 leading-relaxed">
              <p>1. Credentials verified via Bcrypt hash in PostgreSQL.</p>
              <p>2. Go issues signed HMAC-SHA256 JWT token with 24-hour expiration.</p>
              <p>3. On user logout, token JTI is pushed to Redis blacklist with TTL matching remaining token lifespan.</p>
            </div>
          </div>
        </div>
      </div>
    </div>

  {:else if activeTab === 'topology'}
    <!-- Tab 2: Visual Topology Schematic -->
    <div>
      <h2 class="text-base font-semibold text-neutral-900 dark:text-neutral-100 mb-1">Network Topology Schematic</h2>
      <p class="text-xs text-neutral-500 mb-4">Structural visualization of network tiers and isolated boundaries.</p>

      <div class="p-6 sm:p-8 rounded-md border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-[#0A0A0A] font-mono text-xs">
        <!-- Tier 1: Host Ingress -->
        <div class="border border-neutral-300 dark:border-neutral-700 rounded p-4 mb-6 bg-neutral-50 dark:bg-neutral-900/50">
          <div class="flex items-center justify-between mb-3 text-neutral-500 text-xs font-medium">
            <span>Tier 1: Public host ingress boundary</span>
            <span class="text-emerald-600 dark:text-emerald-400 font-semibold">Exposed: Port 1122</span>
          </div>
          <div class="p-3 bg-white dark:bg-[#0A0A0A] border border-neutral-300 dark:border-neutral-700 rounded text-center">
            <div class="text-neutral-900 dark:text-neutral-100 font-bold">Nginx Reverse Proxy (Container: songmoodboard_nginx)</div>
            <div class="text-[11px] text-neutral-500 mt-0.5">Host: 1122 → Container: 80 · SSL Termination · Request Dispatcher</div>
          </div>
        </div>

        <!-- Connection Arrow -->
        <div class="flex justify-center -my-3 text-neutral-400 dark:text-neutral-600">
          │ (Docker Bridge: song-moodboard-net)
        </div>

        <!-- Tier 2: Application Tier -->
        <div class="border border-neutral-300 dark:border-neutral-700 rounded p-4 my-6 bg-neutral-50 dark:bg-neutral-900/50">
          <div class="flex items-center justify-between mb-3 text-neutral-500 text-xs font-medium">
            <span>Tier 2: Application services (no host ports)</span>
            <span class="text-neutral-500">Internal network only</span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="p-3 bg-white dark:bg-[#0A0A0A] border border-neutral-300 dark:border-neutral-700 rounded">
              <div class="text-neutral-900 dark:text-neutral-100 font-semibold">Go REST API (backend:8080)</div>
              <div class="text-[11px] text-neutral-500 mt-1">Static binary · Non-root appuser:1001 · Chi Router</div>
            </div>
            <div class="p-3 bg-white dark:bg-[#0A0A0A] border border-neutral-300 dark:border-neutral-700 rounded">
              <div class="text-neutral-900 dark:text-neutral-100 font-semibold">SvelteKit SSR (frontend:3000)</div>
              <div class="text-[11px] text-neutral-500 mt-1">Bun 1.x runtime · Server-side rendering · TanStack Query</div>
            </div>
          </div>
        </div>

        <!-- Connection Arrow -->
        <div class="flex justify-center -my-3 text-neutral-400 dark:text-neutral-600">
          │ (Zero-Trust Internal IPC)
        </div>

        <!-- Tier 3: Data Tier -->
        <div class="border border-neutral-300 dark:border-neutral-700 rounded p-4 mt-6 bg-neutral-50 dark:bg-neutral-900/50">
          <div class="flex items-center justify-between mb-3 text-neutral-500 text-xs font-medium">
            <span>Tier 3: Data persistence & object storage</span>
            <span class="text-rose-500 font-semibold">Host access blocked</span>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="p-3 bg-white dark:bg-[#0A0A0A] border border-neutral-300 dark:border-neutral-700 rounded">
              <div class="text-neutral-900 dark:text-neutral-100 font-semibold">PostgreSQL 16</div>
              <div class="text-[11px] text-neutral-500 mt-1">postgres:5432 · ACID volume</div>
            </div>
            <div class="p-3 bg-white dark:bg-[#0A0A0A] border border-neutral-300 dark:border-neutral-700 rounded">
              <div class="text-neutral-900 dark:text-neutral-100 font-semibold">Redis 7 Alpine</div>
              <div class="text-[11px] text-neutral-500 mt-1">redis:6379 · Cache & Blacklist</div>
            </div>
            <div class="p-3 bg-white dark:bg-[#0A0A0A] border border-neutral-300 dark:border-neutral-700 rounded">
              <div class="text-neutral-900 dark:text-neutral-100 font-semibold">SeaweedFS S3</div>
              <div class="text-[11px] text-neutral-500 mt-1">seaweedfs:8333 · Media Lake</div>
            </div>
          </div>
        </div>
      </div>
    </div>

  {:else if activeTab === 'matrix'}
    <!-- Tab 3: Container Engineering Matrix Table -->
    <div>
      <h2 class="text-base font-semibold text-neutral-900 dark:text-neutral-100 mb-1">Container Engineering Specification Matrix</h2>
      <p class="text-xs text-neutral-500 mb-4">Complete operational specification of every container running in the Song Moodboard stack.</p>

      <div class="border border-neutral-200 dark:border-neutral-800 rounded-md overflow-x-auto bg-white dark:bg-[#0A0A0A]">
        <table class="w-full text-left text-xs">
          <thead class="bg-neutral-100 dark:bg-neutral-900 text-neutral-600 dark:text-neutral-400 text-xs font-medium border-b border-neutral-200 dark:border-neutral-800">
            <tr>
              <th class="py-3 px-4">Service</th>
              <th class="py-3 px-4">Container Image</th>
              <th class="py-3 px-4">Host Port</th>
              <th class="py-3 px-4">Internal Port</th>
              <th class="py-3 px-4">Process Model</th>
              <th class="py-3 px-4">Isolation Tier</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-200 dark:divide-neutral-800 font-normal">
            {#each serviceSpecs as s}
              <tr class="hover:bg-neutral-50 dark:hover:bg-neutral-900/50 transition-colors">
                <td class="py-3 px-4 font-semibold text-neutral-900 dark:text-neutral-100">{s.name}</td>
                <td class="py-3 px-4 font-mono text-neutral-500">{s.image}</td>
                <td class="py-3 px-4 font-mono {s.hostPort !== 'None (Blocked)' ? 'text-emerald-600 dark:text-emerald-400 font-semibold' : 'text-neutral-400'}">
                  {s.hostPort}
                </td>
                <td class="py-3 px-4 font-mono text-neutral-600 dark:text-neutral-300">{s.internalPort}</td>
                <td class="py-3 px-4 text-neutral-600 dark:text-neutral-400">{s.role}</td>
                <td class="py-3 px-4 font-mono text-neutral-500">{s.isolation}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

  {:else if activeTab === 'cicd'}
    <!-- Tab 4: CI/CD Pipeline Breakdown -->
    <div>
      <h2 class="text-base font-semibold text-neutral-900 dark:text-neutral-100 mb-1">Continuous Integration & Deployment Pipeline</h2>
      <p class="text-xs text-neutral-500 mb-4">Verification gatekeepers defined in .github/workflows/ci-cd.yml executing on push and PR to main.</p>

      <div class="border border-neutral-200 dark:border-neutral-800 rounded-md divide-y divide-neutral-200 dark:divide-neutral-800 bg-white dark:bg-[#0A0A0A] text-xs">
        {#each pipelineStages as st}
          <div class="p-5 flex flex-col md:flex-row md:items-start justify-between gap-4">
            <div class="md:w-1/3">
              <span class="text-xs text-neutral-500 font-medium">Stage {st.id}</span>
              <h3 class="text-sm font-semibold text-neutral-900 dark:text-neutral-100 mt-0.5">{st.title}</h3>
            </div>
            <div class="flex-1">
              <div class="font-mono text-[11px] bg-neutral-100 dark:bg-neutral-900 border border-neutral-200 dark:border-neutral-800 px-3 py-1.5 rounded text-neutral-800 dark:text-neutral-200 mb-2">
                $ {st.command}
              </div>
              <p class="text-neutral-600 dark:text-neutral-400 leading-relaxed">
                {st.desc}
              </p>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>

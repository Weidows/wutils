import React, { useEffect, useState } from 'react';

interface ServiceInfo {
  name: string;
  description: string;
  status: string;
}

type Page = 'dashboard' | 'dsg' | 'ol' | 'buffer' | 'zip' | 'diff' | 'media' | 'gmm' | 'config';

const NAV_ITEMS: { id: Page; label: string; icon: string }[] = [
  { id: 'dashboard', label: 'Dashboard', icon: '📊' },
  { id: 'dsg', label: 'Disk Guard', icon: '💾' },
  { id: 'ol', label: 'Opacity', icon: '👁️' },
  { id: 'buffer', label: 'Buffer', icon: '🗄️' },
  { id: 'zip', label: 'Zip Crack', icon: '🔓' },
  { id: 'diff', label: 'Diff', icon: '📋' },
  { id: 'media', label: 'Media', icon: '🖼️' },
  { id: 'gmm', label: 'GMM', icon: '📦' },
  { id: 'config', label: 'Config', icon: '⚙️' },
];

function App() {
  const [page, setPage] = useState<Page>('dashboard');
  const [services, setServices] = useState<ServiceInfo[]>([]);

  const loadServices = async () => {
    try {
      const svcs = await window.go.main.App.ListServices();
      setServices(svcs);
    } catch (e) {
      console.error('Failed to load services:', e);
    }
  };

  useEffect(() => {
    loadServices();
    const interval = setInterval(loadServices, 2000);
    return () => clearInterval(interval);
  }, []);

  const runningCount = services.filter(s => s.status === 'running').length;

  return (
    <div style={{ display: 'flex', height: '100vh' }}>
      {/* Sidebar */}
      <nav style={{
        width: 220,
        background: 'var(--bg-surface)',
        padding: '1rem 0',
        display: 'flex',
        flexDirection: 'column',
        borderRight: '1px solid rgba(255,255,255,0.05)',
      }}>
        <div style={{
          padding: '0 1rem 1rem',
          borderBottom: '1px solid rgba(255,255,255,0.05)',
          marginBottom: '0.5rem',
        }}>
          <h1 style={{ fontSize: 18, fontWeight: 700, margin: 0 }}>⚡ wutils</h1>
          <span style={{ fontSize: 12, color: 'var(--text-secondary)' }}>
            {runningCount}/{services.length} running
          </span>
        </div>

        {NAV_ITEMS.map(item => (
          <button
            key={item.id}
            onClick={() => setPage(item.id)}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 10,
              padding: '0.6rem 1rem',
              border: 'none',
              background: page === item.id ? 'var(--bg-hover)' : 'transparent',
              color: page === item.id ? 'var(--accent)' : 'var(--text-secondary)',
              cursor: 'pointer',
              fontSize: 14,
              textAlign: 'left',
              width: '100%',
              transition: 'all 0.15s',
            }}
            onMouseEnter={e => { if (page !== item.id) (e.target as HTMLElement).style.background = 'var(--bg-hover)'; }}
            onMouseLeave={e => { if (page !== item.id) (e.target as HTMLElement).style.background = 'transparent'; }}
          >
            <span>{item.icon}</span>
            <span>{item.label}</span>
          </button>
        ))}

        <div style={{ marginTop: 'auto', padding: '1rem', fontSize: 12, color: 'var(--text-secondary)' }}>
          Weidows Utilities v0.1
        </div>
      </nav>

      {/* Main content */}
      <main style={{ flex: 1, padding: '2rem', overflowY: 'auto' }}>
        {page === 'dashboard' && <DashboardPage services={services} onRefresh={loadServices} />}
        {page === 'dsg' && <DSGPage />}
        {page === 'ol' && <OLPage />}
        {page === 'buffer' && <div><h2>Buffer Filesystem</h2><p>Coming soon...</p></div>}
        {page === 'zip' && <ZipPage />}
        {page === 'diff' && <DiffPage />}
        {page === 'media' && <div><h2>Media Grouping</h2><p>Coming soon...</p></div>}
        {page === 'gmm' && <div><h2>Go Mirror Manager</h2><p>Coming soon...</p></div>}
        {page === 'config' && <div><h2>Configuration</h2><p>Coming soon...</p></div>}
      </main>
    </div>
  );
}

// ====== Dashboard ======
function DashboardPage({ services, onRefresh }: { services: ServiceInfo[]; onRefresh: () => void }) {
  const handleStopAll = async () => {
    await window.go.main.App.StopAllServices();
    onRefresh();
  };

  const handleStartAll = async () => {
    await window.go.main.App.StartAllServices();
    onRefresh();
  };

  const toggleService = async (name: string, currentStatus: string) => {
    if (currentStatus === 'running') {
      await window.go.main.App.StopService(name);
    } else {
      await window.go.main.App.StartService(name);
    }
    onRefresh();
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <h2 style={{ margin: 0 }}>Dashboard</h2>
        <div style={{ display: 'flex', gap: 8 }}>
          <button onClick={handleStartAll} style={btnStyle('#22c55e')}>▶ Start All</button>
          <button onClick={handleStopAll} style={btnStyle('#ef4444')}>■ Stop All</button>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: 12 }}>
        {services.map(svc => {
          const isRunning = svc.status === 'running';
          const isError = svc.status === 'error';
          return (
            <div key={svc.name} style={{
              background: 'var(--bg-surface)',
              borderRadius: 8,
              padding: '1rem',
              border: `1px solid ${isError ? 'var(--error)' : 'transparent'}`,
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                <h3 style={{ margin: 0, fontSize: 16 }}>{svc.name}</h3>
                <span style={{
                  width: 10, height: 10, borderRadius: '50%',
                  background: isError ? 'var(--error)' : isRunning ? 'var(--success)' : 'var(--text-secondary)',
                  display: 'inline-block',
                }} />
              </div>
              <p style={{ fontSize: 13, color: 'var(--text-secondary)', margin: '0 0 12px' }}>
                {svc.description}
              </p>
              <button
                onClick={() => toggleService(svc.name, svc.status)}
                style={{
                  ...btnStyle(isRunning ? '#ef4444' : '#22c55e'),
                  fontSize: 12, padding: '4px 12px',
                }}
              >
                {isRunning ? 'Stop' : 'Start'}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// ====== DSG Page ======
function DSGPage() {
  const [status, setStatus] = useState('stopped');

  useEffect(() => {
    const load = async () => setStatus(await window.go.main.App.DSGStatus());
    load();
    const interval = setInterval(load, 2000);
    return () => clearInterval(interval);
  }, []);

  const toggle = async () => {
    if (status === 'running') {
      await window.go.main.App.StopDSG();
    } else {
      await window.go.main.App.StartDSG();
    }
    setStatus(await window.go.main.App.DSGStatus());
  };

  return (
    <div>
      <h2>Disk Sleep Guard</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Prevent external HDDs from sleeping by periodic writes.</p>
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginTop: '1rem' }}>
        <span style={{
          width: 12, height: 12, borderRadius: '50%',
          background: status === 'running' ? 'var(--success)' : 'var(--text-secondary)',
          display: 'inline-block',
        }} />
        <span>Status: <strong>{status}</strong></span>
        <button onClick={toggle} style={btnStyle(status === 'running' ? '#ef4444' : '#22c55e')}>
          {status === 'running' ? 'Stop' : 'Start'}
        </button>
      </div>
    </div>
  );
}

// ====== OL Page ======
function OLPage() {
  const [status, setStatus] = useState('stopped');
  const [windows, setWindows] = useState<any[]>([]);

  useEffect(() => {
    const loadStatus = async () => setStatus(await window.go.main.App.OLStatus());
    loadStatus();
    const interval = setInterval(loadStatus, 2000);
    return () => clearInterval(interval);
  }, []);

  const toggle = async () => {
    if (status === 'running') {
      await window.go.main.App.StopOL();
    } else {
      await window.go.main.App.StartOL();
    }
    setStatus(await window.go.main.App.OLStatus());
  };

  const listWindows = async () => {
    const wins = await window.go.main.App.ListWindows();
    setWindows(wins);
  };

  return (
    <div>
      <h2>Opacity Listener</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Auto-transparency for windows matching rules.</p>
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginTop: '1rem' }}>
        <span style={{
          width: 12, height: 12, borderRadius: '50%',
          background: status === 'running' ? 'var(--success)' : 'var(--text-secondary)',
        }} />
        <span>Status: <strong>{status}</strong></span>
        <button onClick={toggle} style={btnStyle(status === 'running' ? '#ef4444' : '#22c55e')}>
          {status === 'running' ? 'Stop' : 'Start'}
        </button>
        <button onClick={listWindows} style={btnStyle('#4a90d9')}>
          List Windows
        </button>
      </div>

      {windows.length > 0 && (
        <div style={{ marginTop: '1rem', background: 'var(--bg-surface)', borderRadius: 8, padding: '1rem' }}>
          <h3 style={{ margin: '0 0 8px' }}>Visible Windows ({windows.length})</h3>
          <div style={{ maxHeight: 300, overflowY: 'auto', fontSize: 13 }}>
            {windows.map((w, i) => (
              <div key={i} style={{ padding: '4px 0', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                <span style={{ color: 'var(--text-secondary)' }}>{w.handle}</span>
                {' '}{w.title}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

// ====== ZipCrack Page ======
function ZipPage() {
  const [path, setPath] = useState('');
  const [result, setResult] = useState('');
  const [loading, setLoading] = useState(false);

  const crack = async () => {
    if (!path) return;
    setLoading(true);
    setResult('');
    try {
      const pwd = await window.go.main.App.CrackPassword(path);
      setResult(pwd ? `✅ Password found: ${pwd}` : '❌ Password not found');
    } catch (e: any) {
      setResult(`❌ Error: ${e.message || e}`);
    }
    setLoading(false);
  };

  return (
    <div>
      <h2>Zip Password Cracker</h2>
      <div style={{ display: 'flex', gap: 8, marginTop: '1rem' }}>
        <input
          value={path}
          onChange={e => setPath(e.target.value)}
          placeholder="Path to archive (.zip/.7z/.rar)"
          style={{
            flex: 1, padding: '8px 12px', borderRadius: 6, border: '1px solid var(--bg-hover)',
            background: 'var(--bg-surface)', color: 'var(--text-primary)', fontSize: 14,
          }}
        />
        <button onClick={crack} disabled={loading} style={btnStyle('#4a90d9')}>
          {loading ? 'Cracking...' : 'Crack'}
        </button>
      </div>
      {result && (
        <div style={{ marginTop: '1rem', padding: '1rem', background: 'var(--bg-surface)', borderRadius: 8 }}>
          {result}
        </div>
      )}
    </div>
  );
}

// ====== Diff Page ======
function DiffPage() {
  const [inputA, setInputA] = useState('');
  const [inputB, setInputB] = useState('');
  const [result, setResult] = useState<DiffResult | null>(null);

  interface DiffResult {
    missInA: string[];
    missInB: string[];
  }

  const runDiff = async () => {
    if (!inputA || !inputB) return;
    const res = await window.go.main.App.Diff(inputA, inputB);
    setResult(res);
  };

  return (
    <div>
      <h2>File Diff</h2>
      <div style={{ display: 'flex', gap: 8, marginTop: '1rem' }}>
        <input value={inputA} onChange={e => setInputA(e.target.value)} placeholder="inputA.txt" style={inputStyle} />
        <input value={inputB} onChange={e => setInputB(e.target.value)} placeholder="inputB.txt" style={inputStyle} />
        <button onClick={runDiff} style={btnStyle('#4a90d9')}>Compare</button>
      </div>
      {result && (
        <div style={{ marginTop: '1rem', display: 'flex', gap: 16 }}>
          <div style={{ flex: 1, background: 'var(--bg-surface)', borderRadius: 8, padding: '1rem' }}>
            <h3 style={{ margin: '0 0 8px', color: '#f87171' }}>Missing in A ({result.missInA.length})</h3>
            {result.missInA.map((line, i) => <div key={i} style={{ fontSize: 13 }}>{line}</div>)}
          </div>
          <div style={{ flex: 1, background: 'var(--bg-surface)', borderRadius: 8, padding: '1rem' }}>
            <h3 style={{ margin: '0 0 8px', color: '#4ade80' }}>Missing in B ({result.missInB.length})</h3>
            {result.missInB.map((line, i) => <div key={i} style={{ fontSize: 13 }}>{line}</div>)}
          </div>
        </div>
      )}
    </div>
  );
}

// ====== Helpers ======
const btnStyle = (color: string): React.CSSProperties => ({
  background: color,
  color: '#fff',
  border: 'none',
  padding: '8px 16px',
  borderRadius: 6,
  cursor: 'pointer',
  fontSize: 14,
  fontWeight: 600,
  opacity: 0.9,
});

const inputStyle: React.CSSProperties = {
  padding: '8px 12px', borderRadius: 6, border: '1px solid var(--bg-hover)',
  background: 'var(--bg-surface)', color: 'var(--text-primary)', fontSize: 14, flex: 1,
};

export default App;

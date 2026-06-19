import React, { useEffect, useState, useCallback } from 'react';

interface ServiceInfo {
  name: string;
  description: string;
  status: string;
}

interface WindowInfo {
  handle: string;
  title: string;
  opacity: number;
}

interface DiffResult {
  missInA: string[];
  missInB: string[];
}

interface MirrorInfo {
  name: string;
  url: string;
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
  const [toast, setToast] = useState<{ msg: string; type: 'ok' | 'err' } | null>(null);

  const showToast = useCallback((msg: string, type: 'ok' | 'err' = 'ok') => {
    setToast({ msg, type });
    setTimeout(() => setToast(null), 3000);
  }, []);

  const loadServices = useCallback(async () => {
    try {
      const svcs = await window.go.main.App.ListServices();
      setServices(svcs);
    } catch (e) {
      // Silent on first load
    }
  }, []);

  // Listen for tray "navigate" event
  useEffect(() => {
    const unload = window.runtime?.EventsOn('navigate', (target: string) => {
      if (target === 'dashboard') setPage('dashboard');
    });
    loadServices();
    const interval = setInterval(loadServices, 2000);
    return () => {
      if (unload) unload();
      clearInterval(interval);
    };
  }, [loadServices]);

  const runningCount = services.filter(s => s.status === 'running').length;

  return (
    <div style={{ display: 'flex', height: '100vh', position: 'relative' }}>
      {/* Sidebar */}
      <nav style={{
        width: 220, background: 'var(--bg-surface)', padding: '1rem 0',
        display: 'flex', flexDirection: 'column',
        borderRight: '1px solid rgba(255,255,255,0.05)',
      }}>
        <div style={{ padding: '0 1rem 1rem', borderBottom: '1px solid rgba(255,255,255,0.05)', marginBottom: '0.5rem' }}>
          <h1 style={{ fontSize: 18, fontWeight: 700, margin: 0 }}>⚡ wutils</h1>
          <span style={{ fontSize: 12, color: 'var(--text-secondary)' }}>
            {runningCount}/{services.length} running
          </span>
        </div>
        {NAV_ITEMS.map(item => (
          <button key={item.id} onClick={() => setPage(item.id)}
            style={{
              display: 'flex', alignItems: 'center', gap: 10, padding: '0.6rem 1rem',
              border: 'none', background: page === item.id ? 'var(--bg-hover)' : 'transparent',
              color: page === item.id ? 'var(--accent)' : 'var(--text-secondary)',
              cursor: 'pointer', fontSize: 14, textAlign: 'left', width: '100%',
              transition: 'all 0.15s',
            }}
            onMouseEnter={e => { if (page !== item.id) (e.target as HTMLElement).style.background = 'var(--bg-hover)'; }}
            onMouseLeave={e => { if (page !== item.id) (e.target as HTMLElement).style.background = 'transparent'; }}
          >
            <span>{item.icon}</span><span>{item.label}</span>
          </button>
        ))}
        <div style={{ marginTop: 'auto', padding: '1rem', fontSize: 12, color: 'var(--text-secondary)' }}>
          Weidows Utilities v0.1
        </div>
      </nav>

      {/* Main */}
      <main style={{ flex: 1, padding: '2rem', overflowY: 'auto' }}>
        {page === 'dashboard' && <DashboardPage services={services} onRefresh={loadServices} showToast={showToast} />}
        {page === 'dsg' && <DSGPage showToast={showToast} />}
        {page === 'ol' && <OLPage showToast={showToast} />}
        {page === 'buffer' && <BufferPage showToast={showToast} />}
        {page === 'zip' && <ZipPage showToast={showToast} />}
        {page === 'diff' && <DiffPage />}
        {page === 'media' && <MediaPage showToast={showToast} />}
        {page === 'gmm' && <GMMPage showToast={showToast} />}
        {page === 'config' && <ConfigPage showToast={showToast} />}
      </main>

      {/* Toast */}
      {toast && (
        <div style={{
          position: 'fixed', bottom: 24, right: 24, padding: '12px 20px', borderRadius: 8,
          background: toast.type === 'ok' ? '#22c55e' : '#ef4444', color: '#fff',
          fontSize: 14, fontWeight: 600, zIndex: 9999, boxShadow: '0 4px 12px rgba(0,0,0,0.3)',
        }}>
          {toast.msg}
        </div>
      )}
    </div>
  );
}

// ========== Dashboard ==========
function DashboardPage({ services, onRefresh, showToast }: {
  services: ServiceInfo[]; onRefresh: () => void; showToast: (m: string, t?: 'ok' | 'err') => void;
}) {
  const toggleService = async (name: string, currentStatus: string) => {
    try {
      if (currentStatus === 'running') {
        await window.go.main.App.StopService(name);
        showToast(`Stopped ${name}`);
      } else {
        await window.go.main.App.StartService(name);
        showToast(`Started ${name}`);
      }
      onRefresh();
    } catch (e: any) { showToast(e.message || String(e), 'err'); }
  };

  const handleStartAll = async () => {
    const res = await window.go.main.App.StartAllServices();
    showToast(res[0] || 'All started');
    onRefresh();
  };
  const handleStopAll = async () => {
    const res = await window.go.main.App.StopAllServices();
    showToast(res[0] || 'All stopped');
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
              background: 'var(--bg-surface)', borderRadius: 8, padding: '1rem',
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
              <p style={{ fontSize: 13, color: 'var(--text-secondary)', margin: '0 0 12px' }}>{svc.description}</p>
              <button onClick={() => toggleService(svc.name, svc.status)}
                style={{ ...btnStyle(isRunning ? '#ef4444' : '#22c55e'), fontSize: 12, padding: '4px 12px' }}>
                {isRunning ? 'Stop' : 'Start'}
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// ========== DSG ==========
function DSGPage({ showToast }: { showToast: (m: string, t?: 'ok' | 'err') => void }) {
  const [status, setStatus] = useState('stopped');
  useEffect(() => {
    const load = async () => { try { setStatus(await window.go.main.App.DSGStatus()); } catch {} };
    load(); const id = setInterval(load, 2000); return () => clearInterval(id);
  }, []);
  const toggle = async () => {
    try {
      if (status === 'running') { await window.go.main.App.StopDSG(); showToast('DSG stopped'); }
      else { await window.go.main.App.StartDSG(); showToast('DSG started'); }
      setStatus(await window.go.main.App.DSGStatus());
    } catch (e: any) { showToast(e.message || String(e), 'err'); }
  };
  return (
    <div>
      <h2>💾 Disk Sleep Guard</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Prevent external HDDs from sleeping by periodic writes.</p>
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginTop: '1rem' }}>
        <span style={{ width: 12, height: 12, borderRadius: '50%', background: status === 'running' ? 'var(--success)' : 'var(--text-secondary)' }} />
        <span>Status: <strong>{status}</strong></span>
        <button onClick={toggle} style={btnStyle(status === 'running' ? '#ef4444' : '#22c55e')}>
          {status === 'running' ? 'Stop' : 'Start'}
        </button>
      </div>
    </div>
  );
}

// ========== OL ==========
function OLPage({ showToast }: { showToast: (m: string, t?: 'ok' | 'err') => void }) {
  const [status, setStatus] = useState('stopped');
  const [windows, setWindows] = useState<WindowInfo[]>([]);
  useEffect(() => {
    const load = async () => { try { setStatus(await window.go.main.App.OLStatus()); } catch {} };
    load(); const id = setInterval(load, 2000); return () => clearInterval(id);
  }, []);
  const toggle = async () => {
    try {
      if (status === 'running') { await window.go.main.App.StopOL(); showToast('OL stopped'); }
      else { await window.go.main.App.StartOL(); showToast('OL started'); }
      setStatus(await window.go.main.App.OLStatus());
    } catch (e: any) { showToast(e.message || String(e), 'err'); }
  };
  const listWindows = async () => {
    try { setWindows(await window.go.main.App.ListWindows()); } catch (e: any) { showToast(e.message || String(e), 'err'); }
  };
  return (
    <div>
      <h2>👁️ Opacity Listener</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Auto-transparency for windows matching rules.</p>
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginTop: '1rem' }}>
        <span style={{ width: 12, height: 12, borderRadius: '50%', background: status === 'running' ? 'var(--success)' : 'var(--text-secondary)' }} />
        <span>Status: <strong>{status}</strong></span>
        <button onClick={toggle} style={btnStyle(status === 'running' ? '#ef4444' : '#22c55e')}>{status === 'running' ? 'Stop' : 'Start'}</button>
        <button onClick={listWindows} style={btnStyle('#4a90d9')}>List Windows</button>
      </div>
      {windows.length > 0 && (
        <div style={{ marginTop: '1rem', background: 'var(--bg-surface)', borderRadius: 8, padding: '1rem' }}>
          <h3 style={{ margin: '0 0 8px' }}>Visible Windows ({windows.length})</h3>
          <div style={{ maxHeight: 300, overflowY: 'auto', fontSize: 13 }}>
            {windows.map((w, i) => (
              <div key={i} style={{ padding: '4px 0', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                <span style={{ color: 'var(--text-secondary)' }}>{w.handle}</span> {w.title}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

// ========== Buffer ==========
function BufferPage({ showToast }: { showToast: (m: string, t?: 'ok' | 'err') => void }) {
  const [drive, setDrive] = useState('X:');
  const [source, setSource] = useState('');
  const [strategy, setStrategy] = useState('balanced');
  const [mounted, setMounted] = useState(false);

  // Since buffer mount/unmount isn't exposed as a service with persistent status,
  // we track mount state locally.
  const doMount = async () => {
    if (!source) { showToast('Please enter source path', 'err'); return; }
    // Buffer mount is done via the CLI or direct API; for now show guidance
    showToast('Use CLI: wutils buffer mount X: --source <path>');
  };
  const doUnmount = async () => {
    showToast('Use CLI: wutils buffer unmount');
  };

  return (
    <div>
      <h2>🗄️ Buffer Filesystem</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Dokan-based IO buffer virtual filesystem for HDD write optimization.</p>
      <div style={{ marginTop: '1rem', display: 'flex', flexDirection: 'column', gap: 12, maxWidth: 500 }}>
        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <label style={{ minWidth: 80 }}>Drive:</label>
          <select value={drive} onChange={e => setDrive(e.target.value)} style={inputStyle}>
            {['X:', 'Y:', 'Z:', 'W:', 'V:'].map(d => <option key={d} value={d}>{d}</option>)}
          </select>
        </div>
        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <label style={{ minWidth: 80 }}>Source:</label>
          <input value={source} onChange={e => setSource(e.target.value)} placeholder="D:/data" style={{ ...inputStyle, flex: 1 }} />
        </div>
        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <label style={{ minWidth: 80 }}>Strategy:</label>
          <select value={strategy} onChange={e => setStrategy(e.target.value)} style={inputStyle}>
            {['balanced', 'monitoring', 'defrag', 'download', 'migration'].map(s =>
              <option key={s} value={s}>{s}</option>
            )}
          </select>
        </div>
        <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
          <button onClick={doMount} disabled={mounted} style={btnStyle(mounted ? '#666' : '#22c55e')}>Mount</button>
          <button onClick={doUnmount} disabled={!mounted} style={btnStyle(mounted ? '#ef4444' : '#666')}>Unmount</button>
        </div>
        <p style={{ fontSize: 12, color: 'var(--text-secondary)', marginTop: 8 }}>
          Note: Requires Dokan driver installed. The buffer feature uses CLI commands for full control.
        </p>
      </div>
    </div>
  );
}

// ========== ZipCrack ==========
function ZipPage({ showToast }: { showToast: (m: string, t?: 'ok' | 'err') => void }) {
  const [path, setPath] = useState('');
  const [result, setResult] = useState('');
  const [loading, setLoading] = useState(false);
  const crack = async () => {
    if (!path) return;
    setLoading(true); setResult('');
    try {
      const pwd = await window.go.main.App.CrackPassword(path);
      setResult(pwd ? `✅ Password found: ${pwd}` : '❌ Password not found');
      if (pwd) showToast(`Password found: ${pwd}`);
      else showToast('Password not found in dictionary', 'err');
    } catch (e: any) { setResult(`❌ Error: ${e.message || e}`); showToast(e.message || String(e), 'err'); }
    setLoading(false);
  };
  return (
    <div>
      <h2>🔓 Zip Password Cracker</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Crack passwords for encrypted .zip, .7z, .rar files.</p>
      <div style={{ display: 'flex', gap: 8, marginTop: '1rem' }}>
        <input value={path} onChange={e => setPath(e.target.value)}
          placeholder="Path to archive (.zip/.7z/.rar)" style={{ ...inputStyle, flex: 1 }} />
        <button onClick={crack} disabled={loading} style={btnStyle('#4a90d9')}>
          {loading ? 'Cracking... 🔄' : 'Crack'}
        </button>
      </div>
      {result && <div style={{ marginTop: '1rem', padding: '1rem', background: 'var(--bg-surface)', borderRadius: 8 }}>{result}</div>}
    </div>
  );
}

// ========== Diff ==========
function DiffPage() {
  const [inputA, setInputA] = useState('./inputA.txt');
  const [inputB, setInputB] = useState('./inputB.txt');
  const [result, setResult] = useState<DiffResult | null>(null);
  const runDiff = async () => {
    const res = await window.go.main.App.Diff(inputA, inputB);
    setResult(res);
  };
  return (
    <div>
      <h2>📋 File Diff</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Compare two files and find line-level differences (symmetric difference).</p>
      <div style={{ display: 'flex', gap: 8, marginTop: '1rem' }}>
        <input value={inputA} onChange={e => setInputA(e.target.value)} style={inputStyle} />
        <input value={inputB} onChange={e => setInputB(e.target.value)} style={inputStyle} />
        <button onClick={runDiff} style={btnStyle('#4a90d9')}>Compare</button>
      </div>
      {result && (
        <div style={{ marginTop: '1rem', display: 'flex', gap: 16 }}>
          <div style={{ flex: 1, background: 'var(--bg-surface)', borderRadius: 8, padding: '1rem' }}>
            <h3 style={{ margin: '0 0 8px', color: '#f87171' }}>Missing in A ({result.missInA.length})</h3>
            {result.missInA.map((line, i) => <div key={i} style={{ fontSize: 13, padding: '2px 0' }}>{line}</div>)}
            {result.missInA.length === 0 && <span style={{ color: 'var(--text-secondary)', fontSize: 13 }}>None</span>}
          </div>
          <div style={{ flex: 1, background: 'var(--bg-surface)', borderRadius: 8, padding: '1rem' }}>
            <h3 style={{ margin: '0 0 8px', color: '#4ade80' }}>Missing in B ({result.missInB.length})</h3>
            {result.missInB.map((line, i) => <div key={i} style={{ fontSize: 13, padding: '2px 0' }}>{line}</div>)}
            {result.missInB.length === 0 && <span style={{ color: 'var(--text-secondary)', fontSize: 13 }}>None</span>}
          </div>
        </div>
      )}
    </div>
  );
}

// ========== Media ==========
function MediaPage({ showToast }: { showToast: (m: string, t?: 'ok' | 'err') => void }) {
  const [dir, setDir] = useState('');
  const [running, setRunning] = useState(false);
  const doCluster = async () => {
    if (!dir) { showToast('Please enter directory path', 'err'); return; }
    setRunning(true);
    try {
      await window.go.main.App.ClusterMedia(dir);
      showToast('Media clustering complete! Check output/ directory.');
    } catch (e: any) { showToast(e.message || String(e), 'err'); }
    setRunning(false);
  };
  return (
    <div>
      <h2>🖼️ Media Grouping</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Group photos/videos by time (≤12h) and GPS proximity (≤1km).</p>
      <div style={{ display: 'flex', gap: 8, marginTop: '1rem' }}>
        <input value={dir} onChange={e => setDir(e.target.value)}
          placeholder="F:/Pictures/DCIM/Camera" style={{ ...inputStyle, flex: 1 }} />
        <button onClick={doCluster} disabled={running} style={btnStyle('#4a90d9')}>
          {running ? '⏳ Grouping...' : '📂 Group & Copy'}
        </button>
      </div>
      <p style={{ fontSize: 12, color: 'var(--text-secondary)', marginTop: 8 }}>
        Files are grouped into an output/ subdirectory under the input path. Large collections may take time.
      </p>
    </div>
  );
}

// ========== GMM ==========
function GMMPage({ showToast }: { showToast: (m: string, t?: 'ok' | 'err') => void }) {
  const [activeTab, setActiveTab] = useState<'proxy' | 'sumdb'>('proxy');
  const [proxyResult, setProxyResult] = useState('');
  const [sumdbResult, setSumdbResult] = useState('');

  const testSpeed = async () => {
    try { await window.go.main.App.SetGoproxy('default'); showToast('Test via CLI: wutils gmm test'); }
    catch (e: any) { showToast(e.message || String(e), 'err'); }
  };
  const setProxy = async (name: string) => {
    try {
      await window.go.main.App.SetGoproxy(name);
      setProxyResult(`GOPROXY → ${name}`);
      showToast(`GOPROXY set to ${name}`);
    } catch (e: any) { showToast(e.message || String(e), 'err'); }
  };
  const setSumdb = async (name: string) => {
    try {
      await window.go.main.App.SetGosumdb(name);
      setSumdbResult(`GOSUMDB → ${name}`);
      showToast(`GOSUMDB set to ${name}`);
    } catch (e: any) { showToast(e.message || String(e), 'err'); }
  };

  const proxyMirrors = ['goproxy.cn', 'aliyun', 'baidu', 'huawei', 'tencent', 'proxy-io', 'tuna', 'default'];
  const sumdbMirrors = ['google', 'sumdb-io', 'default'];

  return (
    <div>
      <h2>📦 Go Mirror Manager</h2>
      <p style={{ color: 'var(--text-secondary)' }}>Manage GOPROXY and GOSUMDB module proxy mirrors.</p>

      <div style={{ display: 'flex', gap: 8, marginBottom: '1rem', marginTop: '1rem' }}>
        <button onClick={() => setActiveTab('proxy')} style={tabStyle(activeTab === 'proxy')}>GOPROXY</button>
        <button onClick={() => setActiveTab('sumdb')} style={tabStyle(activeTab === 'sumdb')}>GOSUMDB</button>
        <div style={{ flex: 1 }} />
        <button onClick={testSpeed} style={btnStyle('#4a90d9')}>⏱ Test All (CLI)</button>
      </div>

      {activeTab === 'proxy' && (
        <div style={{ background: 'var(--bg-surface)', borderRadius: 8 }}>
          {proxyResult && <div style={{ padding: '8px 16px', color: 'var(--success)', fontSize: 13 }}>{proxyResult}</div>}
          {proxyMirrors.map(m => (
            <div key={m} onClick={() => setProxy(m)}
              style={{ display: 'flex', justifyContent: 'space-between', padding: '10px 16px', cursor: 'pointer', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
              <span>{m}</span>
              <span style={{ color: 'var(--accent)', fontSize: 13 }}>Apply →</span>
            </div>
          ))}
        </div>
      )}

      {activeTab === 'sumdb' && (
        <div style={{ background: 'var(--bg-surface)', borderRadius: 8 }}>
          {sumdbResult && <div style={{ padding: '8px 16px', color: 'var(--success)', fontSize: 13 }}>{sumdbResult}</div>}
          {sumdbMirrors.map(m => (
            <div key={m} onClick={() => setSumdb(m)}
              style={{ display: 'flex', justifyContent: 'space-between', padding: '10px 16px', cursor: 'pointer', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
              <span>{m}</span>
              <span style={{ color: 'var(--accent)', fontSize: 13 }}>Apply →</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

// ========== Config ==========
function ConfigPage({ showToast }: { showToast: (m: string, t?: 'ok' | 'err') => void }) {
  const [yaml, setYaml] = useState('# Loading...');
  const [dirty, setDirty] = useState(false);

  useEffect(() => {
    // Config could be loaded from the backend, for now show placeholder
    setYaml(`# wutils configuration
# Edit this file at ~/.config/wutils/app.yml

app:
  name: wutils
  version: "1.0.0"
  debug: false

logging:
  level: info
  format: json

refresh: 10

cmd:
  dsg:
    parallel: true
    disk:
      - E:
    delay: 30
  ol:
    parallel: true
    delay: 2
    patterns:
      - title: "(XY|xy)plorer"
        opacity: 200
      - title: "设置$"
        opacity: 220
  buffer:
    enable: false
    memory_limit: 67108864
    flush_interval: 10
    strategy: balanced
`);
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setYaml(e.target.value);
    setDirty(true);
  };

  const save = () => {
    showToast('Config saved (simulated — edit ~/.config/wutils/app.yml directly)');
    setDirty(false);
  };

  return (
    <div>
      <h2>⚙️ Configuration</h2>
      <p style={{ color: 'var(--text-secondary)' }}>
        Edit wutils configuration. Changes are applied on next service start or via hot-reload.
      </p>
      <textarea value={yaml} onChange={handleChange}
        style={{
          width: '100%', height: 400, background: 'var(--bg-surface)', color: 'var(--text-primary)',
          border: '1px solid var(--bg-hover)', borderRadius: 8, padding: '1rem', fontSize: 13,
          fontFamily: 'monospace', resize: 'vertical', marginTop: '1rem',
        }} />
      <div style={{ display: 'flex', gap: 8, marginTop: 8 }}>
        <button onClick={save} disabled={!dirty} style={btnStyle(dirty ? '#22c55e' : '#666')}>
          Save Config
        </button>
        <span style={{ fontSize: 12, color: 'var(--text-secondary)', marginLeft: 8, alignSelf: 'center' }}>
          File: ~/.config/wutils/app.yml
        </span>
      </div>
    </div>
  );
}

// ========== Style Helpers ==========
const btnStyle = (color: string): React.CSSProperties => ({
  background: color, color: '#fff', border: 'none', padding: '8px 16px',
  borderRadius: 6, cursor: 'pointer', fontSize: 14, fontWeight: 600,
  opacity: color === '#666' ? 0.5 : 0.9,
});

const inputStyle: React.CSSProperties = {
  padding: '8px 12px', borderRadius: 6, border: '1px solid var(--bg-hover)',
  background: 'var(--bg-surface)', color: 'var(--text-primary)', fontSize: 14,
};

const tabStyle = (active: boolean): React.CSSProperties => ({
  padding: '8px 16px', border: 'none', borderRadius: 6, cursor: 'pointer',
  fontSize: 14, fontWeight: 600,
  background: active ? 'var(--accent)' : 'var(--bg-surface)',
  color: '#fff',
});

export default App;

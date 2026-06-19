/// <reference types="vite/client" />

// Wails runtime injected by the Go backend
declare global {
  interface Window {
    go: {
      main: {
        App: {
          ListServices(): Promise<ServiceInfo[]>;
          StartService(name: string): Promise<void>;
          StopService(name: string): Promise<void>;
          StartAllServices(): Promise<string[]>;
          StopAllServices(): Promise<string[]>;
          StartDSG(): Promise<void>;
          StopDSG(): Promise<void>;
          DSGStatus(): Promise<string>;
          StartOL(): Promise<void>;
          StopOL(): Promise<void>;
          OLStatus(): Promise<string>;
          ListWindows(): Promise<WindowInfo[]>;
          Diff(inputA: string, inputB: string): Promise<DiffResult>;
          CrackPassword(archivePath: string): Promise<string>;
          ClusterMedia(inputDir: string): Promise<string>;
          Extract(mode: string, rootPath: string): Promise<void>;
          SetGoproxy(name: string): Promise<void>;
          SetGosumdb(name: string): Promise<void>;
          Greet(name: string): Promise<string>;
          Now(): Promise<string>;
        };
      };
    };
    runtime?: {
      EventsOn(event: string, callback: (...args: any[]) => void): () => void;
      EventsEmit(event: string, ...args: any[]): void;
      WindowShow(): void;
      WindowHide(): void;
      WindowCenter(): void;
      Quit(): void;
    };
  }
}

interface ServiceInfo {
  name: string;
  description: string;
  status: 'stopped' | 'running' | 'error';
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

export {};

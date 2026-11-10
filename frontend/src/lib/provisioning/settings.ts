import type { PortalSettings, ProviderMode } from "./types";

const SETTINGS_KEY = "forge.settings";

/**
 * Chemin relatif par défaut : en production, l'Ingress route /api vers le BFF
 * sur la même origine que le portail, donc aucun CORS ni secret côté client.
 */
const DEFAULT_API_BASE = "/api/v1";
const DEFAULT_MODE: ProviderMode = "crossplane";

function isMode(value: unknown): value is ProviderMode {
  return value === "mock" || value === "crossplane";
}

export function readSettings(): PortalSettings {
  try {
    const raw = localStorage.getItem(SETTINGS_KEY);
    if (!raw) {
      return { apiBaseUrl: DEFAULT_API_BASE, mode: DEFAULT_MODE };
    }
    const parsed = JSON.parse(raw) as Partial<PortalSettings>;
    return {
      apiBaseUrl: parsed.apiBaseUrl?.trim() || DEFAULT_API_BASE,
      mode: isMode(parsed.mode) ? parsed.mode : DEFAULT_MODE,
    };
  } catch {
    return { apiBaseUrl: DEFAULT_API_BASE, mode: DEFAULT_MODE };
  }
}

export function writeSettings(settings: PortalSettings): void {
  localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings));
}

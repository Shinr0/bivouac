import { readSettings, writeSettings } from "./settings";
import {
  ValidationError,
  type CatalogItem,
  type CreateVmInput,
  type FormOptions,
  type PortalSettings,
  type ProvisioningPort,
  type ProvisioningRequest,
} from "./types";

interface ErrorBody {
  error?: string;
  fields?: Record<string, string>;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const { apiBaseUrl } = readSettings();

  let resp: Response;
  try {
    resp = await fetch(`${apiBaseUrl}${path}`, {
      ...init,
      headers: {
        Accept: "application/json",
        ...(init?.body ? { "Content-Type": "application/json" } : {}),
        ...init?.headers,
      },
    });
  } catch (cause) {
    throw new Error(
      `API Forge injoignable sur ${apiBaseUrl}. Vérifiez Settings ou basculez en mode mock.`,
      { cause },
    );
  }

  if (resp.status === 204) {
    return undefined as T;
  }

  const payload = (await resp.json().catch(() => ({}))) as unknown;

  if (!resp.ok) {
    const body = payload as ErrorBody;
    if (resp.status === 400 && body.fields) {
      throw new ValidationError(body.error ?? "paramètres invalides", body.fields);
    }
    throw new Error(body.error ?? `HTTP ${resp.status}`);
  }

  return payload as T;
}

/**
 * Mode production : le BFF crée un composite Crossplane VirtualMachine dans le
 * cluster. Aucun YAML n'est généré côté navigateur, aucun secret n'est stocké.
 */
export function createCrossplanePort(): ProvisioningPort {
  return {
    async listCatalog() {
      const data = await request<{ items: CatalogItem[] }>("/catalog");
      return data.items;
    },

    async getFormOptions() {
      const data = await request<{ options: FormOptions }>("/catalog");
      return data.options;
    },

    async createRequest(input: CreateVmInput) {
      return request<ProvisioningRequest>("/vms", {
        method: "POST",
        body: JSON.stringify(input),
      });
    },

    async listRequests() {
      const data = await request<{ items: ProvisioningRequest[] }>("/vms");
      return data.items ?? [];
    },

    async deleteRequest(name: string) {
      await request<void>(`/vms/${encodeURIComponent(name)}`, { method: "DELETE" });
    },

    // Le cycle de vie réel est piloté par Crossplane, pas par le portail.
    async advanceMockStatus() {
      return null;
    },

    getSettings(): PortalSettings {
      return readSettings();
    },

    saveSettings(settings: PortalSettings) {
      writeSettings(settings);
    },
  };
}

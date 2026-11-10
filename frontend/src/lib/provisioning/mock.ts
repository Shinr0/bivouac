import { readSettings, writeSettings } from "./settings";
import type {
  CatalogItem,
  CreateVmInput,
  FormOptions,
  PortalSettings,
  ProvisioningPort,
  ProvisioningRequest,
  RequestStatus,
} from "./types";

const REQUESTS_KEY = "forge.requests";

const CATALOG: CatalogItem[] = [
  {
    id: "virtual-machine",
    title: "Virtual Machine",
    description:
      "Instance Scaleway provisionnée par Crossplane via la XRD VirtualMachine.",
    available: true,
  },
  {
    id: "postgres",
    title: "PostgreSQL managé",
    description: "Base de données Scaleway RDB pour les équipes produit.",
    available: false,
  },
  {
    id: "object-bucket",
    title: "Bucket Object Storage",
    description: "Stockage S3 compatible avec cycle de vie et versioning.",
    available: false,
  },
];

/** Reflète les enum de la XRD ; le BFF sert la même liste en mode crossplane. */
const OPTIONS: FormOptions = {
  machineTypes: [
    "DEV1-S",
    "DEV1-M",
    "DEV1-L",
    "PLAY2-PICO",
    "PLAY2-NANO",
    "PLAY2-MICRO",
    "PRO2-XXS",
    "PRO2-XS",
  ],
  zones: ["fr-par-1", "fr-par-2", "nl-ams-1", "nl-ams-2", "pl-waw-1"],
  imageFamilies: ["ubuntu_jammy", "ubuntu_noble", "debian_bookworm", "rockylinux_9"],
  rootVolumeTypes: ["sbs_volume", "l_ssd"],
  minDiskSizeGb: 10,
  maxDiskSizeGb: 1000,
};

const STATUS_FLOW: RequestStatus[] = ["submitted", "provisioning", "ready"];

function readRequests(): ProvisioningRequest[] {
  try {
    const raw = localStorage.getItem(REQUESTS_KEY);
    if (!raw) return [];
    return JSON.parse(raw) as ProvisioningRequest[];
  } catch {
    return [];
  }
}

function writeRequests(requests: ProvisioningRequest[]) {
  localStorage.setItem(REQUESTS_KEY, JSON.stringify(requests));
}

export function createMockPort(): ProvisioningPort {
  return {
    async listCatalog() {
      return CATALOG;
    },

    async getFormOptions() {
      return OPTIONS;
    },

    async createRequest(input: CreateVmInput) {
      const existing = readRequests();
      if (
        existing.some(
          (r) => r.parameters.vmName === input.vmName && r.status !== "failed",
        )
      ) {
        throw new Error(`Une demande active existe déjà pour « ${input.vmName} ».`);
      }
      const request: ProvisioningRequest = {
        id: crypto.randomUUID(),
        serviceId: "virtual-machine",
        parameters: { ...input },
        status: "submitted",
        createdAt: new Date().toISOString(),
      };
      writeRequests([request, ...existing]);
      return request;
    },

    async listRequests() {
      return readRequests();
    },

    async deleteRequest(name: string) {
      writeRequests(readRequests().filter((r) => r.parameters.vmName !== name));
    },

    async advanceMockStatus(id: string) {
      const requests = readRequests();
      const idx = requests.findIndex((r) => r.id === id);
      if (idx < 0) return null;

      const current = requests[idx];
      const at = STATUS_FLOW.indexOf(current.status);
      if (at < 0 || at >= STATUS_FLOW.length - 1) {
        return current;
      }

      const next: ProvisioningRequest = {
        ...current,
        status: STATUS_FLOW[at + 1],
        publicIp: STATUS_FLOW[at + 1] === "ready" ? "51.15.0.42" : undefined,
      };
      requests[idx] = next;
      writeRequests(requests);
      return next;
    },

    getSettings(): PortalSettings {
      return readSettings();
    },

    saveSettings(settings: PortalSettings) {
      writeSettings(settings);
    },
  };
}

export type RequestStatus = "submitted" | "provisioning" | "ready" | "failed";

/**
 * `crossplane` interroge le BFF, qui crée un composite Crossplane dans le
 * cluster. `mock` reste disponible pour travailler sur l'UI sans cluster.
 */
export type ProviderMode = "mock" | "crossplane";

export interface CatalogItem {
  id: string;
  title: string;
  description: string;
  available: boolean;
}

/** Aligné sur spec.parameters de la XRD VirtualMachine (Scaleway). */
export interface VmParameters {
  vmName: string;
  machineType: string;
  zone: string;
  imageFamily: string;
  diskSizeGb: number;
  rootVolumeType: string;
  dataDiskSizeGb: number;
  staticIp: boolean;
}

/** Valeurs acceptées par la XRD, servies par l'API pour éviter de les dupliquer. */
export interface FormOptions {
  machineTypes: string[];
  zones: string[];
  imageFamilies: string[];
  rootVolumeTypes: string[];
  minDiskSizeGb: number;
  maxDiskSizeGb: number;
}

export interface ProvisioningRequest {
  id: string;
  serviceId: string;
  parameters: VmParameters;
  status: RequestStatus;
  message?: string;
  publicIp?: string;
  serverId?: string;
  state?: string;
  requester?: string;
  createdAt: string;
}

export interface PortalSettings {
  /** Base de l'API du BFF. Relative par défaut : même origine que le portail. */
  apiBaseUrl: string;
  mode: ProviderMode;
}

export type CreateVmInput = VmParameters;

/** Erreur de validation renvoyée champ par champ par le BFF. */
export class ValidationError extends Error {
  readonly fields: Record<string, string>;

  constructor(message: string, fields: Record<string, string>) {
    super(message);
    this.name = "ValidationError";
    this.fields = fields;
  }
}

export interface ProvisioningPort {
  listCatalog(): Promise<CatalogItem[]>;
  getFormOptions(): Promise<FormOptions>;
  createRequest(input: CreateVmInput): Promise<ProvisioningRequest>;
  listRequests(): Promise<ProvisioningRequest[]>;
  deleteRequest(name: string): Promise<void>;
  /** Uniquement en mode mock : fait avancer artificiellement le cycle de vie. */
  advanceMockStatus(id: string): Promise<ProvisioningRequest | null>;
  getSettings(): PortalSettings;
  saveSettings(settings: PortalSettings): void;
}

# Contract: ProvisioningPort

```ts
interface ProvisioningPort {
  listCatalog(): Promise<CatalogItem[]>;
  createRequest(input: CreateVmInput): Promise<ProvisioningRequest>;
  listRequests(): Promise<ProvisioningRequest[]>;
  advanceMockStatus?(id: string): Promise<ProvisioningRequest | null>; // mock helper
  getSettings(): PortalSettings;
  saveSettings(settings: PortalSettings): void;
}
```

- `createRequest` in mock mode persists to localStorage key `forge.requests`.
- `createRequest` in github mode POSTs `workflow_dispatch` to `create-vm.yml` then records a local `submitted` request for Status UX.
- Settings key: `forge.settings`.

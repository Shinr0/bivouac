import { createCrossplanePort } from "./crossplane";
import { createMockPort } from "./mock";
import { readSettings } from "./settings";
import type { ProvisioningPort } from "./types";

export type * from "./types";
export { ValidationError } from "./types";

export function getPort(): ProvisioningPort {
  return readSettings().mode === "mock" ? createMockPort() : createCrossplanePort();
}

export type Client = {
  id: string;
  name: string;
  description: string;
  scopes: string[];
  redirect_uris: string[];
  post_logout_redirect_uris: string[];
  grant_types: string[];
  response_types: string[];
  audiences: string[];
  token_endpoint_auth_method: string;
  application_type: string;
  access_token_type: string;
  login_url: string;
  id_token_lifetime_seconds: number;
  dev_mode: boolean;
  clock_skew_seconds: number;
  id_token_userinfo_claims_assertion: boolean;
  created_by: string;
  updated_by: string;
  created_at: string;
  updated_at: string;
};

export type ClientInput = {
  name: string;
  description: string;
  scopes: string[];
  redirect_uris: string[];
  post_logout_redirect_uris: string[];
  grant_types: string[];
  response_types: string[];
  audiences: string[];
  token_endpoint_auth_method: string;
  application_type: string;
  access_token_type: string;
  login_url: string;
  id_token_lifetime_seconds: number;
  dev_mode: boolean;
  clock_skew_seconds: number;
  id_token_userinfo_claims_assertion: boolean;
};

const BASE_URL = "/api/v1/clients";

async function handle<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const body = await res.json().catch(() => ({}) as { error?: string });
    throw new Error(body.error || `Request failed with status ${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export function listClients(page: number, pageSize: number) {
  return fetch(`${BASE_URL}?page=${page}&page_size=${pageSize}`).then((r) =>
    handle<Client[]>(r),
  );
}

export function createClient(input: ClientInput) {
  return fetch(BASE_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  }).then((r) => handle<{ client: Client; secret: string }>(r));
}

export function updateClient(id: string, input: ClientInput) {
  return fetch(`${BASE_URL}/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  }).then((r) => handle<Client>(r));
}

export function deleteClient(id: string) {
  return fetch(`${BASE_URL}/${id}`, { method: "DELETE" }).then((r) => handle<void>(r));
}

"use client";

import { useState } from "react";
import type { FormEvent } from "react";
import {
  Button,
  Checkbox,
  Description,
  Input,
  Label,
  ListBox,
  Modal,
  Select,
  Tabs,
  TextArea,
  TextField,
} from "@heroui/react";
import type { Client, ClientInput } from "@/lib/api";

const APPLICATION_TYPES = ["web", "user_agent", "native"];
const ACCESS_TOKEN_TYPES = ["bearer", "JWT"];
const AUTH_METHODS = ["client_secret_basic", "client_secret_post", "none", "private_key_jwt"];
const GRANT_TYPES = ["authorization_code", "refresh_token", "client_credentials", "implicit"];
const RESPONSE_TYPES = ["code", "id_token", "id_token token"];

const toList = (text: string) =>
  text
    .split(/[\n,]/)
    .map((v) => v.trim())
    .filter(Boolean);

type FormState = {
  name: string;
  description: string;
  redirectUris: string;
  postLogoutRedirectUris: string;
  scopes: string;
  audiences: string;
  grantTypes: string[];
  responseTypes: string[];
  tokenEndpointAuthMethod: string;
  applicationType: string;
  accessTokenType: string;
  loginUrl: string;
  idTokenLifetimeSeconds: string;
  clockSkewSeconds: string;
  devMode: boolean;
  idTokenUserinfoClaimsAssertion: boolean;
};

function initialState(client: Client | null): FormState {
  if (!client) {
    return {
      name: "",
      description: "",
      redirectUris: "",
      postLogoutRedirectUris: "",
      scopes: "openid",
      audiences: "",
      grantTypes: ["authorization_code"],
      responseTypes: ["code"],
      tokenEndpointAuthMethod: "client_secret_basic",
      applicationType: "web",
      accessTokenType: "bearer",
      loginUrl: "",
      idTokenLifetimeSeconds: "3600",
      clockSkewSeconds: "0",
      devMode: false,
      idTokenUserinfoClaimsAssertion: false,
    };
  }
  return {
    name: client.name,
    description: client.description,
    redirectUris: client.redirect_uris.join(", "),
    postLogoutRedirectUris: client.post_logout_redirect_uris.join(", "),
    scopes: client.scopes.join(", "),
    audiences: client.audiences.join(", "),
    grantTypes: client.grant_types,
    responseTypes: client.response_types,
    tokenEndpointAuthMethod: client.token_endpoint_auth_method,
    applicationType: client.application_type,
    accessTokenType: client.access_token_type,
    loginUrl: client.login_url,
    idTokenLifetimeSeconds: String(client.id_token_lifetime_seconds),
    clockSkewSeconds: String(client.clock_skew_seconds),
    devMode: client.dev_mode,
    idTokenUserinfoClaimsAssertion: client.id_token_userinfo_claims_assertion,
  };
}

// key={client-id|"new"} on the parent's usage of this component ensures a
// fresh instance (and fresh form state) whenever a different client is opened.
export function ClientFormModal({
  client,
  onClose,
  onSubmit,
}: {
  client: Client | "new";
  onClose: () => void;
  onSubmit: (input: ClientInput) => Promise<void>;
}) {
  const editing = client !== "new" ? client : null;
  const [form, setForm] = useState<FormState>(() => initialState(editing));
  const [submitting, setSubmitting] = useState(false);

  const update = <K extends keyof FormState>(key: K, value: FormState[K]) => setForm((f) => ({ ...f, [key]: value }));

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    try {
      await onSubmit({
        name: form.name,
        description: form.description,
        redirect_uris: toList(form.redirectUris),
        post_logout_redirect_uris: toList(form.postLogoutRedirectUris),
        scopes: toList(form.scopes),
        audiences: toList(form.audiences),
        grant_types: form.grantTypes,
        response_types: form.responseTypes,
        token_endpoint_auth_method: form.tokenEndpointAuthMethod,
        application_type: form.applicationType,
        access_token_type: form.accessTokenType,
        login_url: form.loginUrl,
        id_token_lifetime_seconds: Number(form.idTokenLifetimeSeconds) || 0,
        clock_skew_seconds: Number(form.clockSkewSeconds) || 0,
        dev_mode: form.devMode,
        id_token_userinfo_claims_assertion: form.idTokenUserinfoClaimsAssertion,
      });
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal.Backdrop isOpen onOpenChange={(open) => !open && onClose()}>
      <Modal.Container placement="auto">
        {/* Fixed (not max-) height so switching tabs — which changes content
            height — doesn't resize the dialog itself, only the body's scroll. */}
        <Modal.Dialog className="flex h-[90vh] flex-col sm:max-w-xl">
          <Modal.CloseTrigger />
          <Modal.Header>
            <Modal.Heading>{editing ? `Edit ${editing.name}` : "New client"}</Modal.Heading>
          </Modal.Header>
          {/* A plain <form> isn't a flex container, so without this it breaks
              the Dialog's flex-col chain and Body can't shrink-to-scroll. */}
          <form onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col">
            <Modal.Body className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto p-6">
              <TextField isRequired name="name" value={form.name} onChange={(v) => update("name", v)}>
                <Label>Name</Label>
                <Input placeholder="My application" />
              </TextField>

              <TextField name="description" value={form.description} onChange={(v) => update("description", v)}>
                <Label>Description</Label>
                <Input placeholder="Optional" />
              </TextField>

              <Tabs>
                <Tabs.ListContainer>
                  <Tabs.List aria-label="Client settings">
                    <Tabs.Tab id="basic">
                      Basic
                      <Tabs.Indicator />
                    </Tabs.Tab>
                    <Tabs.Tab id="advanced">
                      Advanced
                      <Tabs.Indicator />
                    </Tabs.Tab>
                  </Tabs.List>
                </Tabs.ListContainer>

                <Tabs.Panel id="basic" className="flex flex-col gap-4 pt-4">
                  <TextField value={form.redirectUris} onChange={(v) => update("redirectUris", v)}>
                    <Label>Redirect URIs</Label>
                    <TextArea placeholder="https://app.example.com/callback" rows={2} />
                    <Description>Comma or newline separated</Description>
                  </TextField>

                  <TextField value={form.scopes} onChange={(v) => update("scopes", v)}>
                    <Label>Scopes</Label>
                    <Input placeholder="openid, profile, email" />
                  </TextField>

                  <div className="grid grid-cols-2 gap-4">
                    <Select
                      className="w-full"
                      selectionMode="multiple"
                      value={form.grantTypes}
                      onChange={(keys) => update("grantTypes", keys as string[])}
                    >
                      <Label>Grant types</Label>
                      <Select.Trigger>
                        <Select.Value />
                        <Select.Indicator />
                      </Select.Trigger>
                      <Select.Popover>
                        <ListBox selectionMode="multiple">
                          {GRANT_TYPES.map((gt) => (
                            <ListBox.Item key={gt} id={gt} textValue={gt}>
                              {gt}
                              <ListBox.ItemIndicator />
                            </ListBox.Item>
                          ))}
                        </ListBox>
                      </Select.Popover>
                    </Select>

                    <Select
                      className="w-full"
                      selectionMode="multiple"
                      value={form.responseTypes}
                      onChange={(keys) => update("responseTypes", keys as string[])}
                    >
                      <Label>Response types</Label>
                      <Select.Trigger>
                        <Select.Value />
                        <Select.Indicator />
                      </Select.Trigger>
                      <Select.Popover>
                        <ListBox selectionMode="multiple">
                          {RESPONSE_TYPES.map((rt) => (
                            <ListBox.Item key={rt} id={rt} textValue={rt}>
                              {rt}
                              <ListBox.ItemIndicator />
                            </ListBox.Item>
                          ))}
                        </ListBox>
                      </Select.Popover>
                    </Select>
                  </div>

                  <TextField value={form.loginUrl} onChange={(v) => update("loginUrl", v)}>
                    <Label>Login URL template</Label>
                    <Input placeholder="Defaults to the built-in login page" />
                  </TextField>

                  <Checkbox isSelected={form.devMode} onChange={(v) => update("devMode", v)}>
                    <Checkbox.Content>
                      <Checkbox.Control>
                        <Checkbox.Indicator />
                      </Checkbox.Control>
                      Dev mode
                    </Checkbox.Content>
                  </Checkbox>

                  <Checkbox
                    isSelected={form.idTokenUserinfoClaimsAssertion}
                    onChange={(v) => update("idTokenUserinfoClaimsAssertion", v)}
                  >
                    <Checkbox.Content>
                      <Checkbox.Control>
                        <Checkbox.Indicator />
                      </Checkbox.Control>
                      Assert userinfo claims in ID token
                    </Checkbox.Content>
                  </Checkbox>
                </Tabs.Panel>

                <Tabs.Panel id="advanced" className="flex flex-col gap-4 pt-4">
                  <div className="grid grid-cols-2 gap-4">
                    <Select
                      className="w-full"
                      value={form.applicationType}
                      onChange={(key) => update("applicationType", key as string)}
                    >
                      <Label>Application type</Label>
                      <Select.Trigger>
                        <Select.Value />
                        <Select.Indicator />
                      </Select.Trigger>
                      <Select.Popover>
                        <ListBox>
                          {APPLICATION_TYPES.map((v) => (
                            <ListBox.Item key={v} id={v} textValue={v}>
                              {v}
                              <ListBox.ItemIndicator />
                            </ListBox.Item>
                          ))}
                        </ListBox>
                      </Select.Popover>
                    </Select>

                    <Select
                      className="w-full"
                      value={form.accessTokenType}
                      onChange={(key) => update("accessTokenType", key as string)}
                    >
                      <Label>Access token type</Label>
                      <Select.Trigger>
                        <Select.Value />
                        <Select.Indicator />
                      </Select.Trigger>
                      <Select.Popover>
                        <ListBox>
                          {ACCESS_TOKEN_TYPES.map((v) => (
                            <ListBox.Item key={v} id={v} textValue={v}>
                              {v}
                              <ListBox.ItemIndicator />
                            </ListBox.Item>
                          ))}
                        </ListBox>
                      </Select.Popover>
                    </Select>
                  </div>

                  <Select
                    className="w-full"
                    value={form.tokenEndpointAuthMethod}
                    onChange={(key) => update("tokenEndpointAuthMethod", key as string)}
                  >
                    <Label>Token endpoint auth method</Label>
                    <Select.Trigger>
                      <Select.Value />
                      <Select.Indicator />
                    </Select.Trigger>
                    <Select.Popover>
                      <ListBox>
                        {AUTH_METHODS.map((v) => (
                          <ListBox.Item key={v} id={v} textValue={v}>
                            {v}
                            <ListBox.ItemIndicator />
                          </ListBox.Item>
                        ))}
                      </ListBox>
                    </Select.Popover>
                  </Select>

                  <TextField value={form.postLogoutRedirectUris} onChange={(v) => update("postLogoutRedirectUris", v)}>
                    <Label>Post-logout redirect URIs</Label>
                    <TextArea placeholder="https://app.example.com/logout" rows={2} />
                  </TextField>

                  <TextField value={form.audiences} onChange={(v) => update("audiences", v)}>
                    <Label>Audiences</Label>
                    <Input placeholder="Optional" />
                  </TextField>

                  <div className="grid grid-cols-2 gap-4">
                    <TextField
                      type="number"
                      value={form.idTokenLifetimeSeconds}
                      onChange={(v) => update("idTokenLifetimeSeconds", v)}
                    >
                      <Label>ID token lifetime (seconds)</Label>
                      <Input />
                    </TextField>

                    <TextField
                      type="number"
                      value={form.clockSkewSeconds}
                      onChange={(v) => update("clockSkewSeconds", v)}
                    >
                      <Label>Clock skew (seconds)</Label>
                      <Input />
                    </TextField>
                  </div>
                </Tabs.Panel>
              </Tabs>
            </Modal.Body>
            <Modal.Footer>
              <Button type="button" variant="secondary" onPress={onClose}>
                Cancel
              </Button>
              <Button type="submit" isDisabled={submitting}>
                {editing ? "Save changes" : "Create client"}
              </Button>
            </Modal.Footer>
          </form>
        </Modal.Dialog>
      </Modal.Container>
    </Modal.Backdrop>
  );
}

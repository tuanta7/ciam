"use client";

import { useCallback, useEffect, useState } from "react";
import { Button, Chip, Modal, Pagination, Table, toast } from "@heroui/react";
import { ClientFormModal } from "@/components/client-form-modal";
import {
  createClient,
  deleteClient,
  listClients,
  updateClient,
  type Client,
  type ClientInput,
} from "@/lib/api";

const PAGE_SIZE = 10;

function errorMessage(err: unknown, fallback: string) {
  return err instanceof Error ? err.message : fallback;
}

export default function Home() {
  const [clients, setClients] = useState<Client[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState<Client | "new" | null>(null);
  const [newSecret, setNewSecret] = useState<{ name: string; secret: string } | null>(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    try {
      setClients(await listClients(page, PAGE_SIZE));
    } catch (err) {
      toast.danger(errorMessage(err, "Failed to load clients"));
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    // Fetch-on-mount/page-change: syncs the table with the `page` external
    // source, so the loading flag it sets is intentional, not derivable state.
    // eslint-disable-next-line react-hooks/set-state-in-effect
    refresh();
  }, [refresh]);

  const handleSubmit = async (input: ClientInput) => {
    try {
      if (editing && editing !== "new") {
        await updateClient(editing.id, input);
        toast.success("Client updated");
      } else {
        const { client, secret } = await createClient(input);
        toast.success("Client created");
        if (secret) setNewSecret({ name: client.name, secret });
      }
      setEditing(null);
      await refresh();
    } catch (err) {
      toast.danger(errorMessage(err, "Save failed"));
    }
  };

  const handleDelete = async (client: Client) => {
    if (!confirm(`Delete client "${client.name}"? This cannot be undone.`)) return;
    try {
      await deleteClient(client.id);
      toast.success("Client deleted");
      await refresh();
    } catch (err) {
      toast.danger(errorMessage(err, "Delete failed"));
    }
  };

  return (
    <div className="flex-1 p-8 max-w-5xl mx-auto w-full flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">OIDC Clients</h1>
        <Button onPress={() => setEditing("new")}>New client</Button>
      </div>

      <Table>
        <Table.ScrollContainer>
          <Table.Content aria-label="OIDC clients">
            <Table.Header>
              <Table.Column isRowHeader>Name</Table.Column>
              <Table.Column>Application type</Table.Column>
              <Table.Column>Auth method</Table.Column>
              <Table.Column>Grant types</Table.Column>
              <Table.Column>Created</Table.Column>
              <Table.Column>Actions</Table.Column>
            </Table.Header>
            <Table.Body>
              {clients.map((c) => (
                <Table.Row key={c.id}>
                  <Table.Cell>{c.name}</Table.Cell>
                  <Table.Cell>{c.application_type}</Table.Cell>
                  <Table.Cell>{c.token_endpoint_auth_method}</Table.Cell>
                  <Table.Cell>
                    <div className="flex flex-wrap gap-1">
                      {c.grant_types.map((gt) => (
                        <Chip key={gt} size="sm">
                          {gt}
                        </Chip>
                      ))}
                    </div>
                  </Table.Cell>
                  <Table.Cell>{new Date(c.created_at).toLocaleDateString()}</Table.Cell>
                  <Table.Cell>
                    <div className="flex gap-2">
                      <Button size="sm" variant="secondary" onPress={() => setEditing(c)}>
                        Edit
                      </Button>
                      <Button size="sm" variant="danger" onPress={() => handleDelete(c)}>
                        Delete
                      </Button>
                    </div>
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table.Content>
        </Table.ScrollContainer>
      </Table>

      {!loading && clients.length === 0 && (
        <p className="text-center text-sm text-gray-500">No clients yet.</p>
      )}

      <Pagination className="justify-center">
        <Pagination.Content>
          <Pagination.Item>
            <Pagination.Previous isDisabled={page === 1} onPress={() => setPage((p) => p - 1)}>
              <Pagination.PreviousIcon />
              Previous
            </Pagination.Previous>
          </Pagination.Item>
          <Pagination.Item>
            <Pagination.Next
              isDisabled={clients.length < PAGE_SIZE}
              onPress={() => setPage((p) => p + 1)}
            >
              Next
              <Pagination.NextIcon />
            </Pagination.Next>
          </Pagination.Item>
        </Pagination.Content>
      </Pagination>

      {editing !== null && (
        <ClientFormModal
          key={editing === "new" ? "new" : editing.id}
          client={editing}
          onClose={() => setEditing(null)}
          onSubmit={handleSubmit}
        />
      )}

      <Modal.Backdrop
        isOpen={newSecret !== null}
        onOpenChange={(open) => !open && setNewSecret(null)}
      >
        <Modal.Container placement="auto">
          <Modal.Dialog className="sm:max-w-md">
            <Modal.CloseTrigger />
            <Modal.Header>
              <Modal.Heading>Client secret</Modal.Heading>
            </Modal.Header>
            <Modal.Body className="p-6 flex flex-col gap-2">
              <p className="text-sm">
                This is the only time the secret for <strong>{newSecret?.name}</strong> will be
                shown. Copy it now.
              </p>
              <code className="break-all rounded bg-gray-100 dark:bg-gray-800 p-2 text-sm select-all">
                {newSecret?.secret}
              </code>
            </Modal.Body>
            <Modal.Footer>
              <Button slot="close">Done</Button>
            </Modal.Footer>
          </Modal.Dialog>
        </Modal.Container>
      </Modal.Backdrop>
    </div>
  );
}

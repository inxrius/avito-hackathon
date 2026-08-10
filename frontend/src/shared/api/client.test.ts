import * as assert from 'node:assert/strict';
import { test } from 'node:test';

import { APIError, NetworkError, createRecap, getProfiles } from './client.ts';

const recapResponse = {
  schema_version: '2.0',
  id: 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa',
  profile_id: '11111111-1111-1111-1111-111111111111',
  year: 2026,
  profile: { id: '11111111-1111-1111-1111-111111111111', name: 'Марина' },
  generation: {
    algorithm_version: 'test',
    feature_schema_version: 'features-v1',
    activity_hash: 'sha256:01',
    generated_at: '2026-08-08T00:00:00Z',
    narrative: { source: 'template', prompt_version: 'test' },
  },
  theme: { code: 'city', main_district: { code: 'goods', title: 'Товары' } },
  cards: [],
  capabilities: {
    share_available: true,
    explanation_available: true,
    feedback_available: true,
  },
} as const;

function requestURL(input: string | URL | Request): string {
  if (typeof input === 'string') return input;
  if (input instanceof URL) return input.href;
  return input.url;
}

function bodyAsString(body: BodyInit | null | undefined): string {
  if (typeof body !== 'string') {
    throw new Error('expected request body to be a JSON string');
  }
  return body;
}

void test('createRecap sends the backend payload and accepts 201', async () => {
  const originalFetch = globalThis.fetch;
  let capturedURL = '';
  let capturedInit: RequestInit | undefined;

  const fetchMock: typeof fetch = (input, init) => {
    capturedURL = requestURL(input);
    capturedInit = init;
    return Promise.resolve(
      new Response(JSON.stringify(recapResponse), {
        status: 201,
        headers: { 'Content-Type': 'application/json' },
      }),
    );
  };
  globalThis.fetch = fetchMock;

  try {
    const result = await createRecap('11111111-1111-1111-1111-111111111111', 2026);
    assert.equal(result.id, recapResponse.id);
    assert.equal(capturedURL, '/api/v1/recaps');
    assert.equal(capturedInit?.method, 'POST');
    const sentBody: unknown = JSON.parse(bodyAsString(capturedInit?.body));
    assert.deepEqual(sentBody, {
      profile_id: '11111111-1111-1111-1111-111111111111',
      year: 2026,
    });
  } finally {
    globalThis.fetch = originalFetch;
  }
});

void test('createRecap also accepts an existing snapshot with 200', async () => {
  const originalFetch = globalThis.fetch;
  const fetchMock: typeof fetch = () =>
    Promise.resolve(
      new Response(JSON.stringify(recapResponse), {
        status: 200,
        headers: { 'Content-Type': 'application/json' },
      }),
    );
  globalThis.fetch = fetchMock;

  try {
    await assert.doesNotReject(() =>
      createRecap('11111111-1111-1111-1111-111111111111', 2026),
    );
  } finally {
    globalThis.fetch = originalFetch;
  }
});

void test('client maps backend API errors to APIError', async () => {
  const originalFetch = globalThis.fetch;
  const fetchMock: typeof fetch = () =>
    Promise.resolve(
      new Response(
        JSON.stringify({
          code: 'insufficient_activity',
          message: 'not enough activity',
          request_id: 'request-1',
        }),
        { status: 422, headers: { 'Content-Type': 'application/json' } },
      ),
    );
  globalThis.fetch = fetchMock;

  try {
    await assert.rejects(() => getProfiles(), (error: unknown) => {
      assert.ok(error instanceof APIError);
      assert.equal(error.status, 422);
      assert.equal(error.code, 'insufficient_activity');
      assert.equal(error.requestId, 'request-1');
      return true;
    });
  } finally {
    globalThis.fetch = originalFetch;
  }
});

void test('client maps transport failures to NetworkError', async () => {
  const originalFetch = globalThis.fetch;
  const fetchMock: typeof fetch = () => Promise.reject(new TypeError('connection refused'));
  globalThis.fetch = fetchMock;

  try {
    await assert.rejects(() => getProfiles(), NetworkError);
  } finally {
    globalThis.fetch = originalFetch;
  }
});

void test('client maps invalid error bodies to NetworkError', async () => {
  const originalFetch = globalThis.fetch;
  const fetchMock: typeof fetch = () =>
    Promise.resolve(
      new Response('not-json', {
        status: 503,
        headers: { 'Content-Type': 'text/plain' },
      }),
    );
  globalThis.fetch = fetchMock;

  try {
    await assert.rejects(() => getProfiles(), NetworkError);
  } finally {
    globalThis.fetch = originalFetch;
  }
});

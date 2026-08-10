import * as assert from 'node:assert/strict';
import { test } from 'node:test';

import { APIError, NetworkError } from './client.ts';
import { describeFailure } from './errors.ts';

void test('describeFailure treats insufficient activity as a normal non-retryable state', () => {
  const error = new APIError(422, {
    code: 'insufficient_activity',
    message: 'not enough activity',
    request_id: 'request-1',
  });

  const view = describeFailure(error);
  assert.equal(view.retryable, false);
  assert.match(view.title, /мало действий/u);
});

void test('describeFailure marks network errors as retryable', () => {
  const view = describeFailure(new NetworkError(new Error('offline')));

  assert.equal(view.retryable, true);
  assert.match(view.title, /недоступен/u);
});

import mock from 'src/utils/mock';

mock.onAny(/\/eywa\/api\/functions*/).passThrough()
mock.onAny(/\/eywa\/api\/metrics*/).passThrough()
import mock from 'src/utils/mock';

mock.onAny(/\/eywa\/api\/database*/).passThrough();

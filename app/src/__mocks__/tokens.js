import mock from 'src/utils/mock';

mock.onAny(/\/eywa\/api\/tokens*/).passThrough();

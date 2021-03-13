import mock from 'src/utils/mock';

mock.onAny(/\/oauth*/).passThrough();
mock.onAny(/\/users\/*/).passThrough();
mock.onAny(/\/logout\/*/).passThrough();
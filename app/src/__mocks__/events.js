import mock from 'src/utils/mock';

mock.onGet(/\/eywa\/api\/events*/).passThrough();
mock.onPost(/\/eywa\/api\/events*/).passThrough();

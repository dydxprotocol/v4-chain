import generate from '@babel/generator';
import {
  createTypeRegistry,
  createRegistryLoader,
  ServiceMethod
} from './registry';

export const mutations: ServiceMethod[] = [
  {
    methodName: 'joinPool',
    typeUrl: '/cosmos.pools.transfer.v1.MsgJoinPool',
    TypeName: 'MsgJoinPool'
  },
  {
    methodName: 'exitPool',
    typeUrl: '/cosmos.pools.transfer.v1.MsgExitPool',
    TypeName: 'MsgExitPool'
  }
];


const context = {
  addUtil: () => { }
}

const expectCode = (ast) => {
  expect(
    generate(ast).code
  ).toMatchSnapshot();
}
const printCode = (ast) => {
  console.log(
    generate(ast).code
  );
}

it('createTypeRegistry', async () => {
  expectCode(createTypeRegistry(context, mutations));
});

it('createRegistryLoader', async () => {
  expectCode(createRegistryLoader(context));
});
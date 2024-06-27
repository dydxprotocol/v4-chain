import { createAminoConverter } from '../index';
import { ProtoStore, parseProto } from '@osmonauts/proto-parser'
import { camel } from '@osmonauts/utils';
import { prepareContext, expectCode } from '../../../../../test-utils';

const store = new ProtoStore([]);
store.protos = [];
const addRef = ({ filename, content }) => {
    const ref = {
        absolute: filename,
        filename,
        proto: parseProto(content)
    };
    store.protos.push(ref);
};
addRef({
    filename: 'cosmology/example/a.proto',
    content: `
  syntax = "proto3";
  
  package cosmology.finance;
  option go_package = "github.com/cosmology-finance/go";
  
  message MsgTypePackageA {
      enum EnumPackageA {
          NO_HASH = 0;
          SHA256 = 1;
          SHA512 = 2;
          KECCAK = 3;
          RIPEMD160 = 4;
          BITCOIN = 5;
      }
  
      
      enum EnumDuplicateName {
        A = 0;
        B = 1;
        C = 2;
      }
      
      string address = 1;
      EnumPackageA someCoolField    = 2;
      EnumDuplicateName otherField  = 3;
  
    }
  
  `});
addRef({
    filename: 'cosmology/example/b.proto',
    content: `
  syntax = "proto3";
  
  package cosmology.finance;
  option go_package = "github.com/cosmology-finance/go";
  
  import "cosmology/example/a.proto";
  
  message MsgTypePackageB {
      enum EnumTypePackageB {
          NO_HASH = 0;
          SHA256 = 1;
          SHA512 = 2;
          KECCAK = 3;
          RIPEMD160 = 4;
          BITCOIN = 5;
      }
  
      string address = 1;
      EnumTypePackageB myYolo0 = 2;
  
      message AnotherType {
          MsgTypePackageA myType = 3;
      }
  
      AnotherType anotherField = 4;
  
      enum EnumDuplicateName {
        D = 0;
        E = 1;
        F = 2;
      }
  
      EnumDuplicateName otherField  = 5;
  
  }
  
  `});
addRef({
    filename: 'cosmology/example/c.proto',
    content: `
  syntax = "proto3";
  package cosmology.finance;
  option go_package = "github.com/cosmology-finance/go";
  
  import "cosmology/example/b.proto";
  
  message MsgTypePackageC {
      string                                address    = 1;
      cosmology.finance.MsgTypePackageB     awesome    = 2;
  }
  
  service Msg {
      rpc JoinPool(MsgTypePackageC) returns (MsgTypePackageCResponse);
  }
  message MsgTypePackageCResponse {}
  
  `});

store.traverseAll();

describe('cosmology/example/c', () => {
    const {
        context, root, protos
    } = prepareContext(store, 'cosmology/example/c.proto')

    it('AminoConverter', () => {
        context.options.aminoEncoding.casingFn = camel;
        expectCode(createAminoConverter({
            context,
            root,
            name: 'AminoConverter',
            protos
        }))
    })
});

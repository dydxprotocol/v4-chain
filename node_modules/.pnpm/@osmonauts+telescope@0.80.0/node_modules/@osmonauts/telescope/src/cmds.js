
  import _commands_generate from './commands/generate';
import _commands_install from './commands/install';
import _commands_transpile from './commands/transpile';
  const Commands = {};
  Commands['generate'] = _commands_generate;
Commands['install'] = _commands_install;
Commands['transpile'] = _commands_transpile;
  
    export { Commands }; 
  
    
  import _contracts_generate from './contracts/generate';
import _contracts_install from './contracts/install';
import _contracts_message_composer from './contracts/message-composer';
import _contracts_react_query from './contracts/react-query';
import _contracts_recoil from './contracts/recoil';
  const Contracts = {};
  Contracts['generate'] = _contracts_generate;
Contracts['install'] = _contracts_install;
Contracts['message-composer'] = _contracts_message_composer;
Contracts['react-query'] = _contracts_react_query;
Contracts['recoil'] = _contracts_recoil;
  
    export { Contracts }; 
  
    
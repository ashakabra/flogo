import { IConnectorContribution } from "wi-studio/common/models/contrib";
import { WiContributionUtils } from "wi-studio/app/contrib/wi-contrib";
import {Http} from "@angular/http";

export class Connection {
    public name: string;
    public description: string;
    public domain:string;
    public userName: string;
    public password: string;
}

export class Util {
    public static getConnection = (input: any): Connection => {
       let connection: Connection = new Connection();
       for (let configuration of input) {
            switch(configuration.name) { 
                case "name": { 
                  connection.name = configuration.value;
                  break;
                } 
                case "description": { 
                   connection.description = configuration.value;
                   break;
                } 
                case "domain": {
                   connection.domain = configuration.value;
                   break;   
                } 
                case "userName": { 
                   connection.userName = configuration.value;
                   break;  
                }  
                case "password": {
                    connection.password = configuration.value;
                    break;
                }
            }
        }
        return connection;
    }

    public static isDuplicate = (context: IConnectorContribution, configurations: IConnectorContribution[]): boolean => {
        let newCon: Connection =Util.getConnection(context.settings)
        for (let configuration of configurations) {
            let oldCon: Connection =  Util.getConnection(configuration.settings)
            if (oldCon.name === newCon.name && (WiContributionUtils.getUniqueId(configuration) !== WiContributionUtils.getUniqueId(context))) {
                return true;
            }
        }
        return false;
    }
}

// DiplomaContract.js

'use strict';

const { Contract } = require('fabric-contract-api');

class DiplomaContract extends Contract {

    async InitLedger(ctx) {
        const diplomaRecords = [
            {
                studentName:'kunasin techasueb',
                teacherName:'Dom pothingan',
                diplomaNumber:'1',
                subjectTopic:'Fullstack Developement',
                issuer:'CMU-Eleaning',
                issuedate:'2023-3-2',
                begindate:'2023-2-30',
                enddate:'2023-3-2',
            },
            {
                studentName:'kontakan kamfoo',
                teacherName:'kampong woradut',
                diplomaNumber:'2',
                subjectTopic:'Fullstack Developement',
                issuer:'CMU-Eleaning',
                issuedate:'2023-12-15',
                begindate:'2023-12-13',
                enddate:'2023-12-2',
            },
            {
                studentName:'Ronaldo',
                teacherName:'messi',
                diplomaNumber:'3',
                subjectTopic:'How To dribbing',
                issuer:'Barcelona FC',
                issuedate:'2022-3-2',
                begindate:'2022-2-30',
                enddate:'2022-3-2',
            },
            {
                studentName:'stephen curry',
                teacherName:'jame harden',
                diplomaNumber:'4',
                subjectTopic:'How to shoot 3 point',
                issuer:'Golden warriors',
                issuedate:'2021-3-2',
                begindate:'2021-2-30',
                enddate:'2021-3-2',
            },
            {
                studentName:'Buakaw Bunchamek',
                teacherName:'Rodthang jitmueangnon',
                diplomaNumber:'5',
                subjectTopic:'How to knock in first round',
                issuer:'one championship',
                issuedate:'2020-3-2',
                begindate:'2020-2-30',
                enddate:'2020-3-2',
            },
            
        ];

        for (const diplomaRecord of diplomaRecords) {
            // diplomaRecord.docType = 'diplomaRecord';
            // example of how to write to world state deterministically
            // use convetion of alphabetic order
            // we insert data in alphabetic order using 'json-stringify-deterministic' and 'sort-keys-recursive'
            // when retrieving data, in any lang, the order of data will be the same and consequently also the corresonding hash
            await ctx.stub.putState(diplomaRecord.diplomaNumber, Buffer.from(JSON.stringify(diplomaRecord)));
        }
    }
    async createDiploma(ctx, studentName, teacherName, diplomaNumber, subjectTopic, issuer, issuedate, begindate, enddate) {
        // เช็คว่ามีการสร้างใบปริญญานั้นแล้วหรือไม่
        const existingDiploma = await ctx.stub.getState(diplomaNumber);
        if (existingDiploma && existingDiploma.length > 0) {
            throw new Error(`Diploma with number ${diplomaNumber} already exists`);
        }



        // Create a new diploma record
        const diplomaRecord = {
            studentName,
            teacherName,
            diplomaNumber,
            subjectTopic,
            issuer,
            issuedate,
            begindate,
            enddate,
            // Store the public key of the issuer for verification purposes
        };

        // เก็บข้อมูลเข้าในledger
        await ctx.stub.putState(diplomaNumber, Buffer.from(JSON.stringify(diplomaRecord)));

        const transactionID = ctx.stub.getTxID();

        // Return the created diploma information along with the transaction ID
        return `Diploma created successfully. Transaction ID: ${transactionID}`;
    }

    async queryDiploma(ctx, diplomaNumber) {
        // Retrieve the diploma record for the given diplomaNumber
        const diplomaRecord = await ctx.stub.getState(diplomaNumber);

        if (!diplomaRecord || diplomaRecord.length === 0) {
            throw new Error(`Diploma with number ${diplomaNumber} not found`);
        }

        return JSON.parse(diplomaRecord.toString('utf8'));
    }

    async getAllDiplomas(ctx) {
        const iterator = await ctx.stub.getStateByRange('', '');

        const diplomas = [];
        while (true) {
            const result = await iterator.next();

            if (result.value) {
                const diploma = JSON.parse(result.value.value.toString('utf8'));
                diplomas.push(diploma);
            }

            if (result.done) {
                await iterator.close();
                return diplomas;
            }
        }
    }
    async deleteDiploma(ctx, diplomaNumber) {
        // Check if the diploma with the given diplomaNumber exists
        const existingDiploma = await ctx.stub.getState(diplomaNumber);
        if (!existingDiploma || existingDiploma.length === 0) {
            throw new Error(`Diploma with number ${diplomaNumber} not found`);
        }

        // Delete the diploma record from the ledger
        await ctx.stub.deleteState(diplomaNumber);

        const transactionID = ctx.stub.getTxID();

        // Return a message indicating that the diploma was deleted along with the transaction ID
        return `Diploma with number ${diplomaNumber} deleted successfully. Transaction ID: ${transactionID}`;
    }
    // UpdateAsset updates an existing asset in the world state with provided parameters.
    // async UpdateAsset(ctx, diplomaNumber, color, size, owner, appraisedValue) {
    //     const exists = await this.AssetExists(ctx, id);
    //     if (!exists) {
    //         throw new Error(`The asset ${id} does not exist`);
    //     }

    //     // overwriting original asset with new asset
    //     const updatedAsset = {
    //         ID: id,
    //         Color: color,
    //         Size: size,
    //         Owner: owner,
    //         AppraisedValue: appraisedValue,
    //     };
    //     // we insert data in alphabetic order using 'json-stringify-deterministic' and 'sort-keys-recursive'
    //     return ctx.stub.putState(id, Buffer.from(stringify(sortKeysRecursive(updatedAsset))));
    // }

}

module.exports = DiplomaContract;

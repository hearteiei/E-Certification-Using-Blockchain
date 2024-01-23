// DiplomaContract.js

'use strict';

const { Contract } = require('fabric-contract-api');

class DiplomaContract extends Contract {
    async storeDiploma(ctx, studentID, diplomaNumber, subjects, issuerPublicKey) {
        // Check if the transaction submitter is the registration office (issuer)
        const submitterPublicKey = ctx.stub.getCreator().toString('hex');
        if (submitterPublicKey !== issuerPublicKey) {
            throw new Error('Permission denied: Only the registration office can arrange diplomas');
        }

        // Check if diploma with the given diplomaNumber already exists
        const existingDiploma = await ctx.stub.getState(diplomaNumber);
        if (existingDiploma && existingDiploma.length > 0) {
            throw new Error(`Diploma with number ${diplomaNumber} already exists`);
        }

        // Create a new diploma record
        const diplomaRecord = {
            studentID,
            diplomaNumber,
            subjects,
        };

        // Store the diploma record on the ledger
        await ctx.stub.putState(diplomaNumber, Buffer.from(JSON.stringify(diplomaRecord)));
    }

    async queryDiploma(ctx, diplomaNumber) {
        // Retrieve the diploma record for the given diplomaNumber
        const diplomaRecord = await ctx.stub.getState(diplomaNumber);

        if (!diplomaRecord || diplomaRecord.length === 0) {
            throw new Error(`Diploma with number ${diplomaNumber} not found`);
        }

        // Check if the transaction submitter is the student who owns the diploma
        const submitterPublicKey = ctx.stub.getCreator().toString('hex');
        const storedDiploma = JSON.parse(diplomaRecord.toString('utf8'));

        if (storedDiploma.studentID !== submitterPublicKey) {
            throw new Error('Permission denied: You can only view your own diploma information');
        }

        return storedDiploma;
    }

    async getAllDiplomas(ctx) {
        // Get all diploma records from the ledger
        const submitterPublicKey = ctx.stub.getCreator().toString('hex');
        const iterator = await ctx.stub.getStateByRange('', '');
        const diplomas = [];

        for await (const result of iterator) {
            const diploma = JSON.parse(result.value.toString('utf8'));

            // Include only diplomas owned by the current student
            if (diploma.studentID === submitterPublicKey) {
                diplomas.push(diploma);
            }
        }

        return diplomas;
    }
}

module.exports = DiplomaContract;

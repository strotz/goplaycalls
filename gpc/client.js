const Client = {
    _tests: [],
    name: "HTTP Client",
    test(testName, func){
        if (!this._tests) {
            this._tests = [];
        }
        this._tests = this._tests.concat({
            testName: testName,
            func: func,
        });
    },
    assert(condition, message) {
        if (!condition) {
            throw new Error(message);
        }
    },
    log(message){
        // Should log into designated output
        console.log(message);
    },
    exit() {
        // Exit VM
    },
    runTests(response){
        const failures = [];
        if (this._tests) {
            for (const t of this._tests) {
                this.log(`RUN: ${t.testName}`);
                try {
                    t.func(response);
                    this.log(`PASS: ${t.testName}`);
                } catch (e) {
                    this.log(`FAILED: ${t.testName}`);
                    this.log(e.toString());
                    failures.push(e);
                }
            }
        }
        return failures;
    },
    // TODO: globals kept across all sessions/VMs
};

const client = Object.create(Client);
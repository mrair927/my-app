import { expect, test, describe } from 'vitest';
import { exec } from 'child_process';
import { promisify } from 'util';
import fs from 'fs';
import path from 'path';

const execAsync = promisify(exec);

/**
 * @param {any[]} data
 */
async function runPythonFunction(data) {
	const pythonScript = `
import json
import sys
sys.path.append('src/microservices')
from operations import up_or_down

data = json.loads(sys.argv[1])
result = up_or_down(data)
print(result)
    `;

	const { stdout } = await execAsync(`python3 -c "${pythonScript}" '${JSON.stringify(data)}'`);
	return stdout.trim();
}

/**
 * @param {string} filePath
 */
function loadJsonData(filePath) {
	const data = JSON.parse(fs.readFileSync(path.resolve(__dirname, filePath), 'utf8'));
	return data.datapoints;
}

describe('upOrDown function', () => {
	test('should return DOWN for empty data', async () => {
		const result = await runPythonFunction([]);
		expect(result).toBe('DOWN');
	});

	test('should return UP for all up data', async () => {
		const data = Array(100).fill([0, 1234567890]); // 100 data points with all 0s (UP)
		const result = await runPythonFunction(data);
		expect(result).toBe('UP');
	});

	test('should return UP for mixed data with 70% up', async () => {
		const data = [...Array(70).fill([0, 1234567890]), ...Array(30).fill([1, 1234567890])];
		const result = await runPythonFunction(data);
		expect(result).toBe('UP');
	});

	test('should return DOWN for edge case with 60% up', async () => {
		const data = [...Array(60).fill([0, 1234567890]), ...Array(40).fill([1, 1234567890])];
		const result = await runPythonFunction(data);
		expect(result).toBe('DOWN');
	});
});

document.addEventListener('DOMContentLoaded', function () {
    fetchChain();
});

async function fetchChain() {
    try {
        const response = await axios.get('http://localhost:8080/chain');
        const chain = response.data.chain;
        const blockchainDiv = document.getElementById('blockchain');
        blockchainDiv.innerHTML = '';

        // Пропустить первый блок и начать с индекса 1
        for (let i = 1; i < chain.length; i++) {
            const block = chain[i];
            const blockDiv = document.createElement('div');
            blockDiv.className = 'block';
            blockDiv.innerHTML = `
                <h3>Block ${block.index}</h3>
                <p>Timestamp: ${new Date(block.timestamp * 1000).toString()}</p>
                <p>Previous Hash: ${block.previous_hash}</p>
                <p>Proof: ${block.proof}</p>
                <h4>Transactions:</h4>
                <div class="transactions">
                    ${block.transactions.slice(0,-1).map(tx => `
                        <div class="transaction">
                            <p><strong>ID:</strong> ${tx.id}</p>
                            <p>${tx.sender} sent ${tx.amount} to ${tx.recipient}</p>
                        </div>
                    `).join('')}
                </div>
            `;
            blockchainDiv.appendChild(blockDiv);
        }
    } catch (error) {
        console.error('Error fetching chain:', error);
    }
}

async function createTransaction() {
    const sender = document.getElementById('sender').value;
    const recipient = document.getElementById('recipient').value;
    const amount = document.getElementById('amount').value;

    try {
        const response = await axios.post('http://localhost:8080/transactions/new', {
            sender,
            recipient,
            amount: parseInt(amount, 10)
        });
        mineBlock();
        document.getElementById('response').textContent = JSON.stringify(response.data, null, 2);
        fetchChain();
        alertDiv.style.display = 'none';
    } catch (error) {
        console.error('Error creating transaction:', error);
        document.getElementById('response').textContent = 'Error creating transaction';
    }
}

async function mineBlock() {
    try {
        const response = await axios.get('http://localhost:8080/mine');
        document.getElementById('response').textContent = JSON.stringify(response.data, null, 2);
        fetchChain();
    } catch (error) {
        console.error('Error mining block:', error);
        document.getElementById('response').textContent = 'Error mining block';
    }
}

async function getTransactionById() {
    const transactionId = document.getElementById('transactionId').value;
    const responseDiv = document.getElementById('response');

    if (transactionId) {
        try {
            const response = await axios.get(`http://localhost:8080/transactions?id=${transactionId}`);
            const transaction = response.data;
            
            responseDiv.innerHTML = `
                <h3>Transaction Details</h3>
                <p><strong>ID:</strong> ${transaction.id}</p>
                <p><strong>Sender:</strong> ${transaction.sender}</p>
                <p><strong>Recipient:</strong> ${transaction.recipient}</p>
                <p><strong>Amount:</strong> ${transaction.amount}</p>
            `;
        } catch (error) {
            console.error('Error fetching transaction:', error);
            responseDiv.textContent = 'Transaction not found';
        }
    } else {
        responseDiv.textContent = 'Please enter a transaction ID.';
    }
}
function initiateTransaction() {
    const alertDiv = document.getElementById('alert');
    alertDiv.style.display = 'block';
    createTransaction(); // Перенесено из этой функции в mineBlock()
}

function showElement(id) {
    document.getElementById(id).classList.remove('hidden');
}

function hideElement(id) {
    document.getElementById(id).classList.add('hidden');
}

function showError(message) {
    document.getElementById('errorMessage').textContent = message;
    showElement('error');
}

function clearOrderInfo() {
    hideElement('orderInfo');
    hideElement('error');
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU');
}

function formatCurrency(amount, currency) {
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: currency
    }).format(amount / 100);
}

function loadOrder() {
    const orderId = document.getElementById('orderId').value.trim();
    
    if (!orderId) {
        showError('Пожалуйста, введите ID заказа');
        return;
    }

    clearOrderInfo();
    showElement('loading');
    
    fetch(`/orders/${orderId}`)
        .then(response => {
            if (!response.ok) {
                if (response.status === 404) {
                    throw new Error('Заказ не найден');
                } else {
                    throw new Error('Ошибка сервера');
                }
            }
            return response.json();
        })
        .then(order => {
            displayOrder(order);
        })
        .catch(error => {
            showError(error.message);
        })
        .finally(() => {
            hideElement('loading');
        });
}

function displayOrder(order) {
    // Основная информация
    document.getElementById('orderUid').textContent = order.order_uid;
    document.getElementById('trackNumber').textContent = order.track_number;
    document.getElementById('entry').textContent = order.entry;
    document.getElementById('locale').textContent = order.locale;
    document.getElementById('customerId').textContent = order.customer_id;
    document.getElementById('deliveryService').textContent = order.delivery_service;
    document.getElementById('dateCreated').textContent = formatDate(order.date_created);

    // Доставка
    document.getElementById('deliveryName').textContent = order.delivery.name;
    document.getElementById('deliveryPhone').textContent = order.delivery.phone;
    document.getElementById('deliveryZip').textContent = order.delivery.zip;
    document.getElementById('deliveryCity').textContent = order.delivery.city;
    document.getElementById('deliveryAddress').textContent = order.delivery.address;
    document.getElementById('deliveryRegion').textContent = order.delivery.region;
    document.getElementById('deliveryEmail').textContent = order.delivery.email;

    // Платеж
    document.getElementById('paymentTransaction').textContent = order.payment.transaction;
    document.getElementById('paymentCurrency').textContent = order.payment.currency;
    document.getElementById('paymentProvider').textContent = order.payment.provider;
    document.getElementById('paymentAmount').textContent = formatCurrency(order.payment.amount, order.payment.currency);
    document.getElementById('paymentBank').textContent = order.payment.bank;
    document.getElementById('paymentDeliveryCost').textContent = formatCurrency(order.payment.delivery_cost, order.payment.currency);
    document.getElementById('paymentGoodsTotal').textContent = formatCurrency(order.payment.goods_total, order.payment.currency);

    // Товары
    const itemsList = document.getElementById('itemsList');
    itemsList.innerHTML = '';
    
    order.items.forEach((item, index) => {
        const itemCard = document.createElement('div');
        itemCard.className = 'item-card';
        itemCard.innerHTML = `
            <h4>Товар ${index + 1}: ${item.name}</h4>
            <div class="item-details">
                <div><strong>Brand:</strong> ${item.brand}</div>
                <div><strong>Price:</strong> ${formatCurrency(item.price, order.payment.currency)}</div>
                <div><strong>Sale:</strong> ${item.sale}%</div>
                <div><strong>Total Price:</strong> ${formatCurrency(item.total_price, order.payment.currency)}</div>
                <div><strong>Size:</strong> ${item.size}</div>
                <div><strong>Status:</strong> ${item.status}</div>
                <div><strong>Chrt ID:</strong> ${item.chrt_id}</div>
                <div><strong>Nm ID:</strong> ${item.nm_id}</div>
            </div>
        `;
        itemsList.appendChild(itemCard);
    });

    showElement('orderInfo');
}

// Обработка нажатия Enter в поле ввода
document.getElementById('orderId').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        loadOrder();
    }
});

// Автофокус на поле ввода при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('orderId').focus();
});
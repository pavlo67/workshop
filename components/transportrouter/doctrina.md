## Як забезпечити несуперечливість розподіленої доменної системи? І, ширше, загальної системи імен. 

Є варіянти... Основною проблемою видається узгодження конфліктів між автономними доменними-системами, якщо в деякий момент вони вирішують обʼєднатись. 


### Блокчейн-консенсус  

Це доволі затратно по розробці і, головне, не вирішує проблеми конфліктів при обʼєднанні автономних доменних систем. Але в рамках одної системи, яка розвивається з деякого 
неконфліктного стану (наприклад, зі старту з порожнім списком доменів) — блокчейн має переваги перед застосуванням централізованого доменного сервера. До розробки!


### Контроль конфліктів при кожному отриманні пакету.

Надійно, але затратно по рантайм-ресурсах (але можна верифікувати виключно перший пакет з кожної послідовности, вважаючи продовження підтвердженими автоматично).


### Характеризація імен параметром "підтверджене"

Впроваджуємо поняття "підтверджености адреси":
* при прийомі кожного пакета адреса .From вважається "непідтвердженою";     
* якщо є конфлікт, то пакети-відповіді на цю адресу або просто підуть в нікуди, або зафейлять відправку (якщо буде перевірка)


### !!! Контроль IP при отриманнях пакетів

Вважаючи, що транспортні сесії між вузлами системи є тривкими, контролюємо IP-адресу відправника.

Отримуючи кожну нову IP-адресу, перевіряємо, чи цей IP працює в нашій поточній доменній системі і відмовляємо в прийомі, якщо ні (нехай піде спершу і узгодиться!) 

З певною періодичністю повторюємо перевірки вже відомих IP-адрес, якщо з них ідутьпакети.



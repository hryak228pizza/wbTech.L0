package generator

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"github.com/google/uuid"
    m "github.com/hryak228pizza/wbTech.L0/cmd"
)


func NewOrder() m.Order {

	// every order fields
	orderId := uuid.New().String()+"testID"
    track := fmt.Sprintf("%d", rand.Intn(100000))+"testTrack"
    entry := "testEntry"
    locale := "testLocale"
	sign := "testSignature"
    custId := fmt.Sprintf("%d", rand.Intn(1000))+"testCustomer"
    service := "testService"
    shard := strconv.Itoa(rand.Intn(10))
    smid := int64(rand.Intn(100))
	date := time.Now()
    oof := strconv.Itoa(rand.Intn(10))

	// delivery fields
    name := "testName"
    phone := "+123456789"
    zip := "123456"
    city := "testCity"
    addr := "testAddr"
    region := "testRegion"
    email := "testMail@example.com"

    delivery := m.Delivery{
        OrderUID: orderId,
        Name:     &name,
        Phone:    &phone,
        Zip:      &zip,
        City:     &city,
        Address:  &addr,
        Region:   &region,
        Email:    &email,
    }

	// payment fields
    reqId := "testReq"
    currency := "testCurrency"
    provider := "testProvider"
    amount := int64(1 + rand.Intn(9999))
    paymentDT := time.Now().Unix()
    bank := "testBank"
    deliveryCost := int64(1 + rand.Intn(999))
    goodsTotal := amount - deliveryCost
    fee := int64(1 + rand.Intn(99))

    payment := m.Payment{
        Transaction:   orderId,
        RequestID:     &reqId,
        Currency:      &currency,
        Provider:      &provider,
        Amount:        &amount,
        PaymentDT:     &paymentDT,
        Bank:          &bank,
        DeliveryCost:  &deliveryCost,
        GoodsTotal:    &goodsTotal,
        CustomFee:     &fee,
    }

	// generate item structs
    var items []*m.Item
    for i := 0; i < rand.Intn(3)+1; i++ {
        chrt := int64(1 + rand.Intn(999998))
        price := int64(1 + rand.Intn(999))
        rid := fmt.Sprintf("%d", rand.Intn(1000)) + "testRid"
        itemName := fmt.Sprintf("testItem-%d", i+1)
        sale := int64(5 * rand.Intn(11))
        size := strconv.Itoa(rand.Intn(10))
        total := price - (price * sale / 100)
        nm := int64(1 + rand.Intn(999998))
        brand := "testBrand"
        status := int64(202)
        items = append(items, &m.Item{
            ID:          i + 1,
            OrderUID:    orderId,
            ChrtID:      &chrt,
            TrackNumber: &track,
            Price:       &price,
            Rid:         &rid,
            Name:        &itemName,
            Sale:        &sale,
            Size:        &size,
            TotalPrice:  &total,
            NmID:        &nm,
            Brand:       &brand,
            Status:      &status,
        })
    }

    return m.Order{
        OrderUID:          orderId,
        TrackNumber:       &track,
        Entry:             &entry,
        Delivery:          delivery,
        Payment:           payment,
        Items:             items,
        Locale:            &locale,
        InternalSignature: &sign,
        CustomerID:        &custId,
        DeliveryService:   &service,
        ShardKey:          &shard,
        SmID:              &smid,
        DateCreated:       &date,
        OofShard:          &oof,
    }
}



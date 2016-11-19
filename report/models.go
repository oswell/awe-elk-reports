package report

import "time"

type Bill struct {
        BillType               string
        BillingEntity          string
        BillingPeriodEndDate   string
        BillingPeriodStartDate string
        InvoiceId              string
        PayerAccountId         string
}

type Identity struct {
        LineItemId             string
        TimeInterval           string
}

type LineItem struct {
        AvailabilityZone       string
        BlendedCost            float64
        BlendedRate            float64
        CurrencyCode           string
        LineItemDescription    string
        LineItemType           string
        Operation              string
        ProductCode            string
        TaxType                string
        UnblendedCost          float64
        UnblendedRate          float64
        UsageAccountId         string
        UsageAmount            float64
        UsageEndDate           time.Time
        UsageStartDate         time.Time
        UsageType              string
}

type Pricing struct {
        LeaseContractLength    string
        PurchaseOption         string
        Term                   string
}

type Product struct {
        ProductName                  string
        Availability                 string
        ClockSpeed                   string
        CurrentGeneration            string
        DatabaseEngine               string
        DedicatedEbsThroughput       string
        DeploymentOption             string
        Description                  string
        Durability                   string
        EbsOptimized                 bool
        EndpointType                 string
        EngineCode                   string
        EnhancedNetworkingSupported  string
        FeeCode                      string
        FeeDescription               string
        FromLocation                 string
        FromLocationType             string
        Group                        string
        GroupDescription             string
        InstanceFamily               string
        InstanceType                 string
        LicenseModel                 string
        Location                     string
        LocationType                 string
        MaxIopsBurstPerformance      string
        MaxIopsvolume                string
        MaxThroughputvolume          string
        MaxVolumeSize                string
        Memory                       string
        MinVolumeSize                string
        NetworkPerformance           string
        OperatingSystem              string
        Operation                    string
        PhysicalProcessor            string
        PreInstalledSw               string
        ProcessorArchitecture        string
        ProcessorFeatures            string
        ProductFamily                string
        Provisioned                  string
        RequestDescription           string
        RequestType                  string
        RoutingTarget                string
        RoutingType                  string
        Servicecode                  string
        Sku                          string
        Storage                      string
        StorageClass                 string
        StorageMedia                 string
        Tenancy                      string
        ToLocation                   string
        ToLocationType               string
        TransferType                 string
        Usagetype                    string
        Vcpu                         string
        VolumeType                   string
}

type Reservation struct {
        AvailabilityZone string
        ReservationARN   string
}

type Record struct {
        Bill         Bill
        Identity     Identity
        LineItem     LineItem
        Pricing      Pricing
        Product      Product
        Reservation  Reservation
        ResourceTags map[string]string
}

type DetailedBilling struct {
        InvoiceID                  string
        PayerAccountId             string
        LinkedAccountId            string
        RecordType                 string
        RecordId                   string
        ProductName                string
        RateId                     string
        SubscriptionId             string
        PricingPlanId              string
        UsageType                  string
        Operation                  string
        AvailabilityZone           string
        ReservedInstance           string
        ItemDescription            string
        UsageStartDate             string
        UsageEndDate               string
        UsageQuantity              float64
        BlendedRate                float64
        BlendedCost                float64
        UnBlendedRate              float64
        UnBlendedCost              float64
        ResourceId                 string
        ResourceTags               map[string]string
}

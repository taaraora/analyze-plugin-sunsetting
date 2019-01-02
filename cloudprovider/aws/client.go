package aws

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/supergiant/analyze-plugin-sunsetting/cloudprovider"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/pkg/errors"

	"github.com/supergiant/analyze/pkg/plugin/proto"
)

const serviceCode = "AmazonEC2"
// filters
const neededProductFamily = "Compute Instance"
const neededOperatingSystem = "Linux"
const neededPreInstalledSw = "NA"

type Client struct {
	ec2Service        *ec2.EC2
	logger logrus.FieldLogger
	region string
}

//NewClient creates aws client
func NewClient(clientConfig *proto.AwsConfig, logger logrus.FieldLogger) (*Client, error) {
	var region = clientConfig.GetRegion()
	var c = &Client{
		logger:logger,
		region: region,
	}
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load AWS SDK config")
	}

	cfg.Region = region
	c.ec2Service = ec2.New(cfg)

	return c, nil
}

func (c *Client) GetPrices() (map[string][]cloudprovider.ProductPrice, error) {
	var computeInstancesPrices = make(map[string][]cloudprovider.ProductPrice, 0)

	var offeringsURI = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/" + serviceCode + "/current/" + c.region + "/index.json"
	offeringsRaw, err :=  http.Get(offeringsURI)
	if err != nil {
		return nil, errors.Wrap(err, "can't download prices")
	}
	defer offeringsRaw.Body.Close()
	if offeringsRaw.StatusCode != http.StatusOK {
		return nil, errors.Errorf("can't download prices, server returned %s", offeringsRaw.Status)
	}

	// file size is about 40 - 50 megabytes
	offeringBytes, err := ioutil.ReadAll(offeringsRaw.Body)

	type prices struct {
		Products map[string]struct {
			Sku           string `json:"sku"`
			ProductFamily string `json:"productFamily"`
			Attributes    struct {
				InstanceType string `json:"instanceType"`
				Memory       string `json:"memory"`
				Vcpu         string `json:"vcpu"`
				Usagetype    string `json:"usagetype"`
				Tenancy      string `json:"tenancy"`
				OperatingSystem      string `json:"operatingSystem"`
				PreInstalledSw      string `json:"preInstalledSw"`
			} `json:"attributes"`
		} `json:"products"`
		Terms struct {
			// OnDemand map[Sku]map[Sku.offerTermCode]struct{...
			OnDemand map[string]map[string]struct {
				PriceDimensions map[string]struct {
					Unit         string `json:"unit"`
					PricePerUnit struct {
						USDRate string `json:"USD"`
					} `json:"pricePerUnit"`
				} `json:"priceDimensions"`
			} `json:"OnDemand"`
		} `json:"terms"`
	}

	offerings := &prices{}

	err = json.Unmarshal(offeringBytes, offerings)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal prices")
	}

	for productSku, product := range offerings.Products {
		if product.ProductFamily != neededProductFamily ||
			product.Attributes.OperatingSystem != neededOperatingSystem ||
			product.Attributes.PreInstalledSw != neededPreInstalledSw {
			continue
		}

		var newPriceItem = cloudprovider.ProductPrice{
			InstanceType: product.Attributes.InstanceType,
			Memory:       product.Attributes.Memory,
			Vcpu:         product.Attributes.Vcpu,
			Unit:         "",
			ValuePerUnit: "",
			Currency:     "USD",
			UsageType:    product.Attributes.Usagetype,
			Tenancy:      product.Attributes.Tenancy,
		}

		for _, price := range offerings.Terms.OnDemand[productSku] {
			for _, priceDimension := range  price.PriceDimensions {
				newPriceItem.Unit = priceDimension.Unit
				newPriceItem.ValuePerUnit=  priceDimension.PricePerUnit.USDRate
				break
			}
		}

		_, exists := computeInstancesPrices[product.Attributes.InstanceType]
		if !exists {
			computeInstancesPrices[product.Attributes.InstanceType] = make([]cloudprovider.ProductPrice, 0, 0)
		}
		computeInstancesPrices[product.Attributes.InstanceType] = append(computeInstancesPrices[product.Attributes.InstanceType], newPriceItem)
	}

	return computeInstancesPrices, nil
}

func (c *Client) GetComputeInstances() (map[string]cloudprovider.ComputeInstance, error) {
	var instancesRequest = c.ec2Service.DescribeInstancesRequest(nil)
	var result = map[string]cloudprovider.ComputeInstance{}
	describeInstancesResponse, err := instancesRequest.Send()
	if err != nil {
		return nil, err
	}

	for _, instancesReservation := range describeInstancesResponse.Reservations {
		for _, i := range instancesReservation.Instances {
			if i.InstanceId == nil {
				continue
			}

			var instanceType, _ = i.InstanceType.MarshalValue()

			result[*i.InstanceId] = cloudprovider.ComputeInstance{
				InstanceID:   *i.InstanceId,
				InstanceType: instanceType,
			}
		}
	}

	return result, nil
}

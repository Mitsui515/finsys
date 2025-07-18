namespace go finsys.fraud
namespace py fraud.service

struct TransactionData {
    1: string type,
    2: double amount,
    3: string nameOrig,
    4: double oldBalanceOrig,
    5: double newBalanceOrig,
    6: string nameDest,
    7: double oldBalanceDest, 
    8: double newBalanceDest,
    9: optional string timestamp
}

struct FraudPrediction {
    1: bool isFraud,
    2: double fraudProbability,
    3: map<string, double> featureImportance
}

struct EDARequest {
    1: string dataPath,
    2: optional string startDate,
    3: optional string endDate,
    4: optional list<string> focusFeatures
}

struct FraudAnalysisRequest {
    1: TransactionData transaction,
    2: FraudPrediction prediction,
    3: optional list<string> fewShotExamples,
    4: optional string analysisDepth  # "basic", "detailed", "comprehensive"
}

service FraudService {
    FraudPrediction predictFraud(1: TransactionData transaction),
    bool trainModel(1: string dataPath),
    string generateEDAReport(1: EDARequest request),
    string generateFraudAnalysisReport(1: FraudAnalysisRequest request)
}
import twitch.users.email.Test;

public class TestRpc implements Test.EmailBoss {

    @Override
    public Test.UpdateEmailResponse updateEmail(Test.UpdateEmailRequest request) {
        System.out.println("Received updateEmail Request with id: " + request.getUserId());
        return Test.UpdateEmailResponse.newBuilder()
                .setCleanedEmail(request.getNewEmail().replaceAll("[^a-zA-Z_\\-\\.@]", ""))
                .build();
    }
}
